package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
	"tz/internal/models"
	"tz/internal/repository"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type SubscriptionHandler struct {
	repo      *repository.SubscriptionRepository
	validator *validator.Validate
}

func NewSubscriptionHandler(repo *repository.SubscriptionRepository) *SubscriptionHandler {
	return &SubscriptionHandler{
		repo:      repo,
		validator: validator.New()}
}

// @Summary Создать подписку
// @Description Создаёт новую запись о подписке. Все поля обязательны, кроме end_date.
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param request body models.Subscription true "Данные подписки"
// @Success 201 {object} models.Subscription "Подписка успешно создана"
// @Failure 400 {object} map[string]string "Некорректные данные: неверный UUID, дата или формат JSON"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /subscriptions [post]
func (h *SubscriptionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var subsModel models.Subscription

	var sub struct {
		ServiceName string `json:"service_name"`
		Price       int    `json:"price"`
		UserId      string `json:"user_id"`
		StartDate   string `json:"start_date"`
	}

	if err := json.NewDecoder(r.Body).Decode(&sub); err != nil {
		log.Warn().Err(err).Msg("Invalid JSON")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	timeSub, err := time.Parse("01-2006", sub.StartDate)
	if err != nil {
		log.Warn().Err(err).Msg("Validation error")
		http.Error(w, "Validation Error", http.StatusBadRequest)
		return
	}

	if err := h.validator.Struct(sub); err != nil {
		log.Warn().Err(err).Msg("Validation error")
		http.Error(w, "Validation Error", http.StatusBadRequest)
		return
	}
	_, err = uuid.Parse(sub.UserId)
	if err != nil {
		log.Warn().Err(err).Msg("UUID is invalid")
		http.Error(w, "UUID is invalid", http.StatusBadRequest)
		return
	}

	subsModel.Price = sub.Price
	subsModel.ServiceName = sub.ServiceName
	subsModel.UserId = sub.UserId
	subsModel.StartDate = timeSub

	created, err := h.repo.Create(r.Context(), &subsModel)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create subscription")
		http.Error(w, "Internal server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

// @Summary Получить подписку по ID
// @Description Возвращает данные о подписке по её уникальному идентификатору
// @Tags subscriptions
// @Produce json
// @Param id path int true "Subscription ID"
// @Success 200 {object} models.Subscription "Данные подписки"
// @Failure 400 {object} map[string]string "Некорректный ID"
// @Failure 404 {object} map[string]string "Подписка не найдена"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /subscriptions/{id} [get]
func (h *SubscriptionHandler) Retrieve(w http.ResponseWriter, r *http.Request) {
	var sub models.Subscription
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	sub.ID = id
	if err != nil {
		log.Error().Err(err)
		http.Error(w, "", http.StatusBadRequest)
	}

	obj, err := h.repo.Retrieve(r.Context(), &sub)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create subscription")
		http.Error(w, "Internal server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusFound)
	json.NewEncoder(w).Encode(obj)
}

// @Summary Обновить подписку
// @Description Полностью обновляет данные существующей подписки по ID
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path int true "Subscription ID"
// @Param request body models.Subscription true "Новые данные подписки"
// @Success 200 {object} models.Subscription "Подписка успешно обновлена"
// @Failure 400 {object} map[string]string "Некорректные данные (неверный UUID, дата, ID)"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /subscriptions/{id} [put]
func (h *SubscriptionHandler) Update(w http.ResponseWriter, r *http.Request) {
	var sub models.Subscription
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	sub.ID = id
	if err != nil {
		log.Error().Err(err)
		http.Error(w, "", http.StatusBadRequest)
	}

	if err := json.NewDecoder(r.Body).Decode(&sub); err != nil {
		log.Warn().Err(err).Msg("Invalid JSON")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if err := h.validator.Struct(sub); err != nil {
		log.Warn().Err(err).Msg("Validation error")
		http.Error(w, "Validation Error", http.StatusBadRequest)
		return
	}

	_, err = uuid.Parse(sub.UserId)
	if err != nil {
		log.Warn().Err(err).Msg("UUID is invalid")
		http.Error(w, "UUID is invalid", http.StatusBadRequest)
		return
	}

	updated, err := h.repo.Update(r.Context(), &sub)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create subscription")
		http.Error(w, "Internal server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updated)
}

// @Summary Удалить подписку
// @Description Удаляет подписку по её ID
// @Tags subscriptions
// @Param id path int true "Subscription ID"
// @Success 200 {object} map[string]string "Подписка успешно удалена"
// @Failure 400 {object} map[string]string "Некорректный ID"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /subscriptions/{id} [delete]
func (h *SubscriptionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	var sub models.Subscription
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	sub.ID = id
	if err != nil {
		log.Error().Err(err)
		http.Error(w, "", http.StatusBadRequest)
	}

	err = h.repo.Delete(r.Context(), &sub)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create subscription")
		http.Error(w, "Internal server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

// @Summary Список всех подписок
// @Description Возвращает список всех подписок
// @Tags subscriptions
// @Produce json
// @Success 200 {array} models.Subscription "Список подписок"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /subscriptions [get]
func (h *SubscriptionHandler) List(w http.ResponseWriter, r *http.Request) {
	var subs []models.Subscription

	list, err := h.repo.List(r.Context(), &subs)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create subscription")
		http.Error(w, "Internal server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusFound)
	json.NewEncoder(w).Encode(list)
}

// @Summary Рассчитать общую стоимость подписок
// @Description Возвращает сумму стоимости всех активных подписок за указанный период с фильтрацией
// @Tags subscriptions
// @Produce json
// @Param start_period query string true "Начало периода (MM-YYYY)"
// @Param user_id query string false "ID пользователя (UUID)"
// @Param service_name query string false "Название сервиса"\
// @Failure 400 {string} string "start_period обязателен"
// @Failure 500 {string} string "Внутренняя ошибка"
// @Router /subscriptions/total [get]
func (h *SubscriptionHandler) Total(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	startPeriod := params.Get("start_period")
	endPeriod := params.Get("end_period")

	userID := params.Get("user_id")
	serviceName := params.Get("service_name")

	total, err := h.repo.CalculateTotal(r.Context(), startPeriod, endPeriod, userID, serviceName)
	if err != nil {
		log.Error().Err(err).Msg("Failed to calculate total")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(total)
}
