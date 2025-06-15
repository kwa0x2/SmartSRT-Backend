package route

import (
	"github.com/gin-gonic/gin"
	"github.com/kwa0x2/AutoSRT-Backend/api/http/delivery"
	"github.com/kwa0x2/AutoSRT-Backend/config"
	"github.com/kwa0x2/AutoSRT-Backend/domain"
	"github.com/kwa0x2/AutoSRT-Backend/repository"
	"github.com/kwa0x2/AutoSRT-Backend/usecase"
	"github.com/resend/resend-go/v2"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func NewContactRoute(env *config.Env, group *gin.RouterGroup, db *mongo.Database, resendClient *resend.Client) {
	rr := repository.NewResendRepository(resendClient)

	cd := &delivery.ContactDelivery{
		ContactUseCase: usecase.NewContactUseCase(repository.NewBaseRepository[*domain.Contact](db)),
		ResendUseCase:  usecase.NewResendUseCase(rr),
		Env:            env,
	}

	contactRoute := group.Group("/contact")
	{
		contactRoute.POST("", cd.Create)
	}

}
