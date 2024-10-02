package handler

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/amicie-monami/music-library/internal/domain/dto"
	"github.com/amicie-monami/music-library/pkg/httpkit"
	"github.com/pawpawchat/core/pkg/response"
)

type songDataGetter interface {
	GetSongs(aggregation map[string]any) ([]dto.SongWithDetails, error)
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
// @Success 200 {array} []dto.SongWithDetails "Список песен, прошедших аггрегацию данных."
// @Failure 400 {object} dto.Error "Неверный запрос, некорректные значения параметров."
// @Failure 500 {object} dto.Error "Внутреняя ошибка сервера."
func GetSongsData(repo songDataGetter) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params, err := parseGetSongsDataQueryParams(r)
		if err != nil {
			slog.Info(err.Error())
			response.Json().BadRequest().Body(body{"error": err.Error()}).MustWrite(w)
			return
		}

		songsWithDetails, err := repo.GetSongs(params)
		if err != nil {
			slog.Info(err.Error())
			// needs refactoring: to get rid of the "magic" error
			response.Json().InternalError().Body(dto.Error{Message: "Internal server error"})
			return
		}

		response.Json().OK().Body(body{"songs": songsWithDetails}).MustWrite(w)
		slog.Info("200")
	})
}

func parseGetSongsDataQueryParams(r *http.Request) (map[string]any, error) {
	filterMap, err := parseGetSongsDataFilterParams(r)
	if err != nil {
		return nil, err
	}

	limit, offset, err := parsePaginationParams(r)
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

	availableValues := map[string]string{
		"song_id":      "song_id",
		"group":        "group_name",
		"song":         "song_name",
		"release_date": "release_date",
		"link":         "link",
		"text":         "text",
	}

	fields := strings.Split(strings.TrimSpace(fieldsParam), " ")
	for idx := range fields {
		if _, ok := availableValues[fields[idx]]; !ok {
			return "", fmt.Errorf("invalid value in fields param: %s", fields[idx])
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

	availableValues := map[string]struct{}{
		"song_id":      {},
		"song_name":    {},
		"groups":       {},
		"link":         {},
		"release_date": {},
		"text":         {},
	}

	filterMap := make(map[string]any)
	params := strings.Split(filter, ",")

	for idx := range params {
		param := strings.Split(params[idx], "=")

		if param[0] == "song_id" {
			// add validate for song_id parameter
		}

		if len(param) != 2 {
			return nil, fmt.Errorf("invalid param %s", params[idx])
		}

		if _, ok := availableValues[param[0]]; !ok {
			return nil, fmt.Errorf("invalid value in filter param: %s", params[idx])
		}

		filterMap[param[0]] = param[1]
	}

	return filterMap, nil
}
