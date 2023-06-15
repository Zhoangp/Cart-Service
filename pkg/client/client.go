package client

import (
	"fmt"
	"github.com/Zhoangp/Cart-Service/config"
	"github.com/Zhoangp/Cart-Service/pb/course"
	"google.golang.org/grpc"
)

func InitCourseServiceClient(c *config.Config) (course.CourseServiceClient, error) {
	// using WithInsecure() because no SSL running
	cc, err := grpc.Dial(c.OtherServices.CourseUrl, grpc.WithInsecure())
	if err != nil {
		fmt.Println("Could not connect:", err)
		return nil, err
	}
	return course.NewCourseServiceClient(cc), nil
}
