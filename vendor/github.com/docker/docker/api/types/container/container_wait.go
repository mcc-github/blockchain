package container










type ContainerWaitOKBodyError struct {

	
	Message string `json:"Message,omitempty"`
}



type ContainerWaitOKBody struct {

	
	
	Error *ContainerWaitOKBodyError `json:"Error"`

	
	
	StatusCode int64 `json:"StatusCode"`
}
