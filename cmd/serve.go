package cmd

import (
	"github.com/soprasteria/docktor/server"
	"github.com/soprasteria/docktor/server/models"
	"github.com/soprasteria/docktor/server/modules/email"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// ServeCmd represents the serve command
var ServeCmd = &cobra.Command{
	Use:   "serve",
	Short: "Launch Docktor server",
	Long:  `Docktor server will listen on 0.0.0.0:8080`,
	Run: func(cmd *cobra.Command, args []string) {
		email.InitSMTPConfiguration()
		models.Connect()
		server.New()
	},
}

func init() {

	// Get configuration from command line flags
	ServeCmd.Flags().StringP("mongo-addr", "m", "localhost:27017", "URL to access MongoDB")
	ServeCmd.Flags().StringP("mongo-username", "", "", "A user which has access to MongoDB")
	ServeCmd.Flags().StringP("mongo-password", "", "", "Password of the mongo user")
	ServeCmd.Flags().StringP("redis-addr", "", "localhost:6379", "URL to access Redis")
	ServeCmd.Flags().StringP("redis-password", "", "", "Redis password. Optional")
	ServeCmd.Flags().StringP("jwt-secret", "j", "dev-docktor-secret", "Secret key used for JWT token authentication. Change it in your instance")
	ServeCmd.Flags().StringP("reset-pwd-secret", "", "dev-docktor-reset-pwd-to-change", "Secret key used when resetting the password. Change it in your instance")
	ServeCmd.Flags().StringP("bcrypt-pepper", "p", "dev-docktor-bcrypt", "Pepper used in password generation. Change it in your instance")
	ServeCmd.Flags().StringP("env", "e", "prod", "dev or prod")
	ServeCmd.Flags().String("ldap-address", "", "LDAP full address like : ldap.server:389. Optional")
	ServeCmd.Flags().String("ldap-baseDN", "", "BaseDN. Optional")
	ServeCmd.Flags().String("ldap-domain", "", "Domain of the user. Optional")
	ServeCmd.Flags().String("ldap-bindDN", "", "DN of system account. Optional")
	ServeCmd.Flags().String("ldap-bindPassword", "", "Password of system account. Optional")
	ServeCmd.Flags().String("ldap-searchFilter", "", "LDAP request to find users. Optional")
	ServeCmd.Flags().String("ldap-attr-username", "cn", "LDAP attribute for username of users.")
	ServeCmd.Flags().String("ldap-attr-firstname", "givenName", "LDAP attribute for firstname of users.")
	ServeCmd.Flags().String("ldap-attr-lastname", "sn", "LDAP attribute for lastname of users.")
	ServeCmd.Flags().String("ldap-attr-realname", "cn", "LDAP attribute for firstname of users.")
	ServeCmd.Flags().String("ldap-attr-email", "mail", "LDAP attribute for lastname of users.")
	ServeCmd.Flags().String("smtp-server", "", "SMTP server with its port.")
	ServeCmd.Flags().String("smtp-user", "", "SMTP user for authentication.")
	ServeCmd.Flags().String("smtp-password", "", "SMTP password for authentication.")
	ServeCmd.Flags().String("smtp-sender", "", "Email used as sender of emails")
	ServeCmd.Flags().String("smtp-identity", "", "Identity of the sender")
	ServeCmd.Flags().String("smtp-logo", "", "Link or data URI to a logo image that will be used in header of emails sent by Docktor. Default is Docktor logo.")
	ServeCmd.Flags().String("engine-transitiontimeout", "2h", "Timeout in duration that will cancel an engine transition when reached")

	// Bind env variables.
	_ = viper.BindPFlag("server.mongo.addr", ServeCmd.Flags().Lookup("mongo-addr"))
	_ = viper.BindPFlag("server.mongo.username", ServeCmd.Flags().Lookup("mongo-username"))
	_ = viper.BindPFlag("server.mongo.password", ServeCmd.Flags().Lookup("mongo-password"))
	_ = viper.BindPFlag("server.redis.addr", ServeCmd.Flags().Lookup("redis-addr"))
	_ = viper.BindPFlag("server.redis.password", ServeCmd.Flags().Lookup("redis-password"))
	_ = viper.BindPFlag("auth.jwt-secret", ServeCmd.Flags().Lookup("jwt-secret"))
	_ = viper.BindPFlag("auth.reset-pwd-secret", ServeCmd.Flags().Lookup("reset-pwd-secret"))
	_ = viper.BindPFlag("auth.bcrypt-pepper", ServeCmd.Flags().Lookup("bcrypt-pepper"))
	_ = viper.BindPFlag("ldap.address", ServeCmd.Flags().Lookup("ldap-address"))
	_ = viper.BindPFlag("ldap.baseDN", ServeCmd.Flags().Lookup("ldap-baseDN"))
	_ = viper.BindPFlag("ldap.domain", ServeCmd.Flags().Lookup("ldap-domain"))
	_ = viper.BindPFlag("ldap.bindDN", ServeCmd.Flags().Lookup("ldap-bindDN"))
	_ = viper.BindPFlag("ldap.bindPassword", ServeCmd.Flags().Lookup("ldap-bindPassword"))
	_ = viper.BindPFlag("ldap.searchFilter", ServeCmd.Flags().Lookup("ldap-searchFilter"))
	_ = viper.BindPFlag("ldap.attr.username", ServeCmd.Flags().Lookup("ldap-attr-username"))
	_ = viper.BindPFlag("ldap.attr.firstname", ServeCmd.Flags().Lookup("ldap-attr-firstname"))
	_ = viper.BindPFlag("ldap.attr.lastname", ServeCmd.Flags().Lookup("ldap-attr-lastname"))
	_ = viper.BindPFlag("ldap.attr.realname", ServeCmd.Flags().Lookup("ldap-attr-realname"))
	_ = viper.BindPFlag("ldap.attr.email", ServeCmd.Flags().Lookup("ldap-attr-email"))
	_ = viper.BindPFlag("smtp.server", ServeCmd.Flags().Lookup("smtp-server"))
	_ = viper.BindPFlag("smtp.user", ServeCmd.Flags().Lookup("smtp-user"))
	_ = viper.BindPFlag("smtp.password", ServeCmd.Flags().Lookup("smtp-password"))
	_ = viper.BindPFlag("smtp.sender", ServeCmd.Flags().Lookup("smtp-sender"))
	_ = viper.BindPFlag("smtp.identity", ServeCmd.Flags().Lookup("smtp-identity"))
	_ = viper.BindPFlag("smtp.logo", ServeCmd.Flags().Lookup("smtp-logo"))
	_ = viper.BindPFlag("engine.transitiontimeout", ServeCmd.Flags().Lookup("engine-transitiontimeout"))
	_ = viper.BindPFlag("env", ServeCmd.Flags().Lookup("env"))
	RootCmd.AddCommand(ServeCmd)

}
