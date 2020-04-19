package model

// TestModel used for test
type TestModel struct {
	Identity    string
	TestParam   string `bson:"test_param_2"`
	EmptyString string
	EmptyInt64  int64
}
