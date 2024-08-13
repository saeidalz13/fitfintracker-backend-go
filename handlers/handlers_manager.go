package handlers

type HandlersManager struct {
	AuthHandler    *AuthHandlersManager
	FinanceHandler *FinanceHandlersManager
	FitnessHandler *FitnessHandlersManager
}

func NewHandlersConfig(
	ahm *AuthHandlersManager,
	finhm *FinanceHandlersManager,
	fithm *FitnessHandlersManager,
) *HandlersManager {
	return &HandlersManager{
		AuthHandler:    ahm,
		FinanceHandler: finhm,
		FitnessHandler: fithm,
	}
}
