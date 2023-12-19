package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type UserMetrics struct {
	createdUsersCounter          prometheus.Counter
	increaseBalanceAmountCounter prometheus.Counter
}

func NewUserMetrics() *UserMetrics {
	return &UserMetrics{
		createdUsersCounter: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: "mtc_homework",
			Name:      "created_users",
			Help:      "Number of created users",
		}),
		increaseBalanceAmountCounter: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: "mtc_homework",
			Name:      "increase_balance_amount",
			Help:      "Total amount of money for increase balance operations",
		}),
	}
}

func (userMetrics *UserMetrics) IncCreatedUsersCounter() {
	userMetrics.createdUsersCounter.Inc()
}

func (userMetrics *UserMetrics) IncIncreaseBalanceCounter(amount float64) {
	userMetrics.increaseBalanceAmountCounter.Add(amount)
}
