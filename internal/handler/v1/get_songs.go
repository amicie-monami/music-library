package handler

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/amicie-monami/music-library/internal/domain/dto"
	"github.com/amicie-monami/music-library/pkg/httpkit"
)

type songDataGetter interface {
	GetSongs(ctx context.Context, aggregation map[string]any) ([]*dto.SongWithDetails, error)
}

// @Summary Получение данных библиотеки
// @Description Метод позволяет получить данные библиотеки, поддерживает пагинацию и фильтрацию по всем полям.
// @Version 0.0.1
// @Router /songs [get]
// @Tags Songs
// @Accept json
// @Produce json
// @Param limit query string false "Количество песен, которое необходимо верунть. Стандартное значение 10, предельное 1000."
// @Param offset query string false "Смещение, необходимое для выборки определенного подмножества песен. Стандартное значение 0."
// @Param fields query string false "Список полей, которые необходимо вернуть. Допустимые значения: [song_id, group, song, release_date, link, text]. Зачения передаются через знак ”+”,например: fields=song_id+release_date."
// @Param filter query string false "Фильтр, с помощью которого происходит аггрегация данных. Допустимые значения: [song_id, song_name, groups, release_date, link, text]. Значения передаются через знак ”,”, например: filter=song_id=1,groups=нервы+жщ. Описание каждого параметра приведено ниже."
// @Param (filter)song_id query string false "Параметр описывает фильтр для идентификатора песни. Поддерживает равенство на одно значение и выборку с помощью операторов сравнения: gt(>), ge(>=), le(<=), lt(<). Пример: filter=song_id=gt+2+lt+8."
// @Param (filter)groups query string false "Параметр описывает фильтр для названия группы. Названия групп передаются через знак ”+”.Чувствителен к регистру, пробелы в названиях заменяются знаком ”_”. Пример: filter=groups=Noize_MC+мы."
// @Param (filter)song_name query string false "Параметр описывает фильтр для названий песен. Поддерживает оператор * регулярных выражений, нечувствителен к регистру, множественные значения передаются через знак ”+”. Пример: filter=song_name=Lil\*+\*eva\*."
// @Param (filter)release_date query string false "Параметр описывает фильтр для даты релиза песни. Поддерживает прямое равенство, операторы сравнения (см. (filter)song_id) и установку границ с помощью знака ”-”. Пример: filter=release_date=01.01.2023-05.05.2024 (start_date-end_date)."
// @Param (filter)link query string false "Параметр описывает фильтр для ссылки на песню. Поддерживает оператор * регулярных выражений, нечувствителен к регистру, множественные значения передаются чреез знак ”+”. Пример: filter=link=\*yandex\*+\*spotify\*."
// @Param (filter)text query string false "Параметр описывает фильтр для текста песни. Поддерживает оператор * регулярных выражений, нечувствителен к регистру, множественные значения передаются чреез знак ”+”. Пример: filter=text=\*батюшка\*+\*ленин\*."
// @Success 200 {array} dto.GetSongsResponse "Список песен, прошедших аггрегацию данных."
// @Failure 400 {object} dto.Error "Неверный запрос, некорректные значения параметров."
// @Failure 500 {object} dto.Error "Внутреняя ошибка сервера."
func GetSongs(repo songDataGetter) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params, err := parseGetSongsDataQueryParams(r)
		if err != nil {
			sendError(w, err)
			return
		}

		songs, err := repo.GetSongs(r.Context(), params)
		if err != nil {
			sendError(w, err)
			return
		}

		slog.Info("songs have been successfully filtered", "count", len(songs))
		responseBody := dto.GetSongsResponse{Songs: songs}
		httpkit.Ok(w, responseBody)
	})
}

func parseGetSongsDataQueryParams(r *http.Request) (map[string]any, error) {
	filterMap, err := parseGetSongsDataFilterParams(r)
	if err != nil {
		return nil, err
	}

	limit, err := parseLimitParam(r)
	if err != nil {
		return nil, err
	}

	offset, err := parseOffsetParam(r)
	if err != nil {
		return nil, err
	}

	fields, err := parseGetSongsFieldsParam(r)
	if err != nil {
		return nil, err
	}

	return map[string]any{
		"filter": filterMap,
		"limit":  limit,
		"offset": offset,
		"fields": fields,
	}, nil
}

func parseGetSongsFieldsParam(r *http.Request) (string, error) {
	fieldsParam := httpkit.GetStrParam("fields", r)
	if fieldsParam == "" {
		return "", nil
	}

	decodedFieldsParam, err := url.QueryUnescape(fieldsParam)
	if err != nil {
		return "", dto.NewError(400, "failed to decode url", "parseGetSongsFieldsParam", nil, nil)
	}

	availableValues := map[string]string{
		"song_id":      "song_id",
		"group":        "group_name",
		"song":         "song_name",
		"release_date": "release_date",
		"link":         "link",
		"text":         "text",
	}

	fieldsParam = strings.ReplaceAll(decodedFieldsParam, "+", " ")
	fields := strings.Split(strings.TrimSpace(fieldsParam), " ")

	for idx := range fields {
		if _, ok := availableValues[fields[idx]]; !ok {
			return "", dto.NewError(400, "unknown parameter", "parseGetSongsFieldsParam", fields[idx], nil)
		}
		//replace query field names with the database column names
		fieldsParam = strings.ReplaceAll(fieldsParam, fields[idx], availableValues[fields[idx]])
	}

	return fieldsParam, nil
}

func parseGetSongsDataFilterParams(r *http.Request) (map[string]any, error) {
	filter := httpkit.GetStrParam("filter", r)
	if filter == "" {
		return nil, nil
	}

	decodedFilter, err := url.QueryUnescape(filter)
	if err != nil {
		return nil, dto.NewError(400, "failed to decode url", "parseGetSongsFieldsParam", nil, nil)
	}

	availableValues := map[string]struct{}{
		"song_id":      {},
		"song_name":    {},
		"groups":       {},
		"link":         {},
		"release_date": {},
		"text":         {},
	}

	params := strings.Split(strings.ReplaceAll(decodedFilter, "+", " "), ",")

	filterMap := make(map[string]any)
	for idx := range params {
		param := strings.Split(params[idx], "=")

		if param[0] == "song_id" {
			// add validate for song_id parameter
		}

		if len(param) != 2 || param[1] == "" {
			return nil, dto.NewError(400, "missing value for filter", "parseGetSongsDataFilterParams", params[idx], nil)
		}

		if _, ok := availableValues[param[0]]; !ok {
			return nil, dto.NewError(400, "invalid filter key", "parseGetSongsDataFilterParams", param[0], nil)
		}

		filterMap[param[0]] = param[1]
	}

	return filterMap, nil
}

func parseOffsetParam(r *http.Request) (int64, error) {
	offsetParam := r.URL.Query().Get("offset")
	if offsetParam == "" {
		return 0, nil
	}

	//check if the offset param is valid
	offset, err := checkOffsetParamValue(offsetParam)
	if err != nil {
		return 0, err
	}

	return offset, nil
}

func checkOffsetParamValue(offsetParam string) (int64, error) {
	offset, err := strconv.ParseInt(offsetParam, 10, 64)
	if err != nil {
		details := fmt.Sprintf("offset=%s", offsetParam)
		return 0, dto.NewError(400, "offset param must be a num", "parseOffsetParam", details, nil)
	}

	if offset < 0 {
		details := fmt.Sprintf("offset=%d, but param must be >= 0", offset)
		return 0, dto.NewError(400, "invalid offset param", "checkOffsetParamValue", details, nil)
	}

	return offset, nil
}

func parseLimitParam(r *http.Request) (int64, error) {
	limitParam := r.URL.Query().Get("limit")
	if limitParam == "" {
		return 10, nil
	}

	//check if the limit param is valid
	limit, err := checkLimitParamValue(limitParam)
	if err != nil {
		return 0, err
	}

	return limit, nil
}

func checkLimitParamValue(limitParam string) (int64, error) {
	//parse limit param
	limit, err := strconv.ParseInt(limitParam, 10, 64)
	if err != nil {
		details := fmt.Sprintf("limit=%s, but it must be >= 0", limitParam)
		return 0, dto.NewError(400, "invalid limit param", "checkLimitParamValue", details, nil)
	}
	// if limit == 0 {
	// 	return 10, nil
	// }

	if limit > 1000 {
		return 1000, nil
	}

	if limit < 1 {
		details := fmt.Sprintf("limit=%d, but must be >= 1", limit)
		return 0, dto.NewError(400, "invalid limit param", "checkLimitParamValue", details, nil)
	}

	return limit, nil
}
