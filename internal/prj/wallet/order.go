package wallet

import "arch_course/internal/prj"

func (o *Order) GetSteps() (steps map[uint8]Step) {
	if o.Type == TypeBuyOrder {
		return map[uint8]Step{
			prj.SagaTypeExchanger: createExchangerOrder,
			prj.SagaTypeBcgateway: createBcgatewayOrder,
		}
	}

	return map[uint8]Step{
		prj.SagaTypeExchanger: createExchangerOrder,
	}
}
