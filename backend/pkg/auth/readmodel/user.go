package readmodel

type UserReadModel struct {
	UserKey  string `gorm:"index:idx_user_read_model_uq,unique;primaryKey"`
	Email    string
	Username string
	Version  int `gorm:"index:idx_user_read_model_uq,unique"`
}
