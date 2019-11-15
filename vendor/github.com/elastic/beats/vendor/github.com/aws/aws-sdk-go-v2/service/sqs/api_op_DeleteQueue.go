// Code generated by private/model/cli/gen-api/main.go. DO NOT EDIT.

package sqs

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/internal/awsutil"
	"github.com/aws/aws-sdk-go-v2/private/protocol"
	"github.com/aws/aws-sdk-go-v2/private/protocol/query"
)

// Please also see https://docs.aws.amazon.com/goto/WebAPI/sqs-2012-11-05/DeleteQueueRequest
type DeleteQueueInput struct {
	_ struct{} `type:"structure"`

	// The URL of the Amazon SQS queue to delete.
	//
	// Queue URLs and names are case-sensitive.
	//
	// QueueUrl is a required field
	QueueUrl *string `type:"string" required:"true"`
}

// String returns the string representation
func (s DeleteQueueInput) String() string {
	return awsutil.Prettify(s)
}

// Validate inspects the fields of the type to determine if they are valid.
func (s *DeleteQueueInput) Validate() error {
	invalidParams := aws.ErrInvalidParams{Context: "DeleteQueueInput"}

	if s.QueueUrl == nil {
		invalidParams.Add(aws.NewErrParamRequired("QueueUrl"))
	}

	if invalidParams.Len() > 0 {
		return invalidParams
	}
	return nil
}

// Please also see https://docs.aws.amazon.com/goto/WebAPI/sqs-2012-11-05/DeleteQueueOutput
type DeleteQueueOutput struct {
	_ struct{} `type:"structure"`
}

// String returns the string representation
func (s DeleteQueueOutput) String() string {
	return awsutil.Prettify(s)
}

const opDeleteQueue = "DeleteQueue"

// DeleteQueueRequest returns a request value for making API operation for
// Amazon Simple Queue Service.
//
// Deletes the queue specified by the QueueUrl, regardless of the queue's contents.
// If the specified queue doesn't exist, Amazon SQS returns a successful response.
//
// Be careful with the DeleteQueue action: When you delete a queue, any messages
// in the queue are no longer available.
//
// When you delete a queue, the deletion process takes up to 60 seconds. Requests
// you send involving that queue during the 60 seconds might succeed. For example,
// a SendMessage request might succeed, but after 60 seconds the queue and the
// message you sent no longer exist.
//
// When you delete a queue, you must wait at least 60 seconds before creating
// a queue with the same name.
//
// Cross-account permissions don't apply to this action. For more information,
// see see Grant Cross-Account Permissions to a Role and a User Name (http://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/sqs-customer-managed-policy-examples.html#grant-cross-account-permissions-to-role-and-user-name)
// in the Amazon Simple Queue Service Developer Guide.
//
//    // Example sending a request using DeleteQueueRequest.
//    req := client.DeleteQueueRequest(params)
//    resp, err := req.Send(context.TODO())
//    if err == nil {
//        fmt.Println(resp)
//    }
//
// Please also see https://docs.aws.amazon.com/goto/WebAPI/sqs-2012-11-05/DeleteQueue
func (c *Client) DeleteQueueRequest(input *DeleteQueueInput) DeleteQueueRequest {
	op := &aws.Operation{
		Name:       opDeleteQueue,
		HTTPMethod: "POST",
		HTTPPath:   "/",
	}

	if input == nil {
		input = &DeleteQueueInput{}
	}

	req := c.newRequest(op, input, &DeleteQueueOutput{})
	req.Handlers.Unmarshal.Remove(query.UnmarshalHandler)
	req.Handlers.Unmarshal.PushBackNamed(protocol.UnmarshalDiscardBodyHandler)
	return DeleteQueueRequest{Request: req, Input: input, Copy: c.DeleteQueueRequest}
}

// DeleteQueueRequest is the request type for the
// DeleteQueue API operation.
type DeleteQueueRequest struct {
	*aws.Request
	Input *DeleteQueueInput
	Copy  func(*DeleteQueueInput) DeleteQueueRequest
}

// Send marshals and sends the DeleteQueue API request.
func (r DeleteQueueRequest) Send(ctx context.Context) (*DeleteQueueResponse, error) {
	r.Request.SetContext(ctx)
	err := r.Request.Send()
	if err != nil {
		return nil, err
	}

	resp := &DeleteQueueResponse{
		DeleteQueueOutput: r.Request.Data.(*DeleteQueueOutput),
		response:          &aws.Response{Request: r.Request},
	}

	return resp, nil
}

// DeleteQueueResponse is the response type for the
// DeleteQueue API operation.
type DeleteQueueResponse struct {
	*DeleteQueueOutput

	response *aws.Response
}

// SDKResponseMetdata returns the response metadata for the
// DeleteQueue request.
func (r *DeleteQueueResponse) SDKResponseMetdata() *aws.Response {
	return r.response
}
