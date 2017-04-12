package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/mkideal/cli"
)

var (
	version string
)

type opt struct {
	Help     bool   `cli:"h,help" usage:"display help"`
	Version  bool   `cli:"version" usage:"display version and revision"`
	Prefix   string `cli:"prefix" usage:"specify prefix of key"`
	RmPrefix bool   `cli:"rm-prefix" usage:"remove prefix of key"`
	Region   string `cli:"region" usage:"specify aws region"`
}

func main() {
	var prefix string
	var region string
	var rmPrefix bool

	cli.Run(&opt{}, func(ctx *cli.Context) error {
		argv := ctx.Argv().(*opt)
		if argv.Help {
			ctx.String(ctx.Usage())
			os.Exit(1)
		}
		if argv.Version {
			ctx.String(fmt.Sprintf("%s\n", version))
			os.Exit(1)
		}
		prefix = argv.Prefix
		region = argv.Region
		rmPrefix = argv.RmPrefix

		return nil
	})

	// setup aws client
	config := aws.NewConfig()
	if region != "" {
		config.WithRegion(region)
	}
	svc := ssm.New(session.New(), config)

	// find key list from parameter store
	keys, err1 := findParameterKeys(svc, prefix)
	if err1 != nil {
		fmt.Println(err1.Error())
		os.Exit(1)
	} else if len(keys) == 0 {
		fmt.Println("no keys")
		os.Exit(1)
	}

	// get parameters
	pz, err2 := getParametersWithKeys(svc, keys)
	if err2 != nil {
		fmt.Println(err2.Error())
		os.Exit(1)
	}

	// output stdout k=v
	output(pz, prefix, rmPrefix)
}

func findParameterKeys(svc *ssm.SSM, prefix string) ([]string, error) {
	// TODO: nexttoken loop
	params := &ssm.DescribeParametersInput{}
	resp, err := svc.DescribeParameters(params)
	if err != nil {
		return nil, err
	}
	keys := make([]string, 0)
	for _, p := range resp.Parameters {
		if prefix != "" && strings.Index(*p.Name, prefix) < 0 {
			continue
		}
		keys = append(keys, *p.Name)
	}
	return keys, nil
}

func getParametersWithKeys(svc *ssm.SSM, keys []string) ([]*ssm.Parameter, error) {
	awsKeys := make([]*string, 0)
	for _, k := range keys {
		awsKeys = append(awsKeys, aws.String(k))
	}
	params := &ssm.GetParametersInput{
		Names:          awsKeys,
		WithDecryption: aws.Bool(true),
	}
	resp, err := svc.GetParameters(params)
	if err != nil {
		return nil, err
	}
	return resp.Parameters, nil
}

func output(pz []*ssm.Parameter, prefix string, rmPrefix bool) {
	for _, p := range pz {
		k := convertKeyToEnv(removeKeyPrefix(*p.Name, prefix, rmPrefix))
		v := *p.Value
		fmt.Printf("export %s=%s\n", k, v)
	}
}

func removeKeyPrefix(key string, prefix string, rmPrefix bool) string {
	if rmPrefix {
		return strings.Replace(key, prefix, "", 1)
	} else {
		return key
	}
}

func convertKeyToEnv(key string) string {
	k1 := strings.Replace(key, ".", "_", -1)
	k2 := strings.Replace(k1, "/", "_", -1)
	k3 := strings.Replace(k2, "-", "_", -1)
	return strings.ToUpper(k3)
}
