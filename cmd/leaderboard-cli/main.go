package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/bigflood/leaderboard/pkg/http_client"
	"github.com/spf13/cobra"
)

func newClient(cmd *cobra.Command) (*http_client.Client, error) {
	endpoint, err := cmd.Flags().GetString("endpoint")
	if err != nil {
		return nil, err
	}
	return http_client.New(endpoint), nil
}

func init() {
	rootCmd.PersistentFlags().StringP("endpoint", "e", "http://localhost:8080", "endpoint (required)")
	rootCmd.AddCommand(userCountCmd)
	rootCmd.AddCommand(setUserCmd)
	rootCmd.AddCommand(getUserCmd)
	rootCmd.AddCommand(getRanksCmd)
}

var rootCmd = &cobra.Command{
	Use:           "leaderboard-cli",
	Short:         "leaderboard api client",
	SilenceUsage:  true,
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

var userCountCmd = &cobra.Command{
	Use: "usercount",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		client, err := newClient(cmd)
		if err != nil {
			return err
		}

		count, err := client.UserCount(ctx)
		if err != nil {
			return err
		}

		fmt.Println(count)
		return nil
	},
}

var setUserCmd = &cobra.Command{
	Use: "setuser [flags] userId score",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return errors.New("invalid number of arguments")
		}

		userId := args[0]
		score, err := strconv.Atoi(args[1])
		if err != nil {
			return err
		}

		ctx := context.Background()

		client, err := newClient(cmd)
		if err != nil {
			return err
		}

		if err := client.SetUser(ctx, userId, score); err != nil {
			return err
		}

		return nil
	},
}

var getUserCmd = &cobra.Command{
	Use: "getuser [flags] userId",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("invalid number of arguments")
		}

		userId := args[0]

		ctx := context.Background()

		client, err := newClient(cmd)
		if err != nil {
			return err
		}

		user, err := client.GetUser(ctx, userId)
		if err != nil {
			return err
		}

		fmt.Printf("%+v\n", user)
		return nil
	},
}

var getRanksCmd = &cobra.Command{
	Use: "getranks [flags] rank count",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return errors.New("invalid number of arguments")
		}

		rank, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}

		count, err := strconv.Atoi(args[1])
		if err != nil {
			return err
		}

		ctx := context.Background()

		client, err := newClient(cmd)
		if err != nil {
			return err
		}

		users, err := client.GetRanks(ctx, rank, count)
		if err != nil {
			return err
		}

		for _, user := range users {
			fmt.Printf("%+v\n", user)
		}
		return nil
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
