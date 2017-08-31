/*
Package types contains all available entities that compose Docktor application.
These entities are structures persited to database and exposed within a REST-like API.

Daemon and Site

Docker is used in Docktor to deploy web tools automatically. A daemon defines an entrypoint to a Docker daemon. A site defines where the daemon is located for statistical purposes.
This daemon can be opened or securized with TLS configuration.

User

A user is someone who is registered and authorized to use Docktor. The person can be either administrator or simple user. An administrator is able to do and see everything. Rights of simple user are defined in groups to whom they belong to.

Catalog

A catalog defines available services and templates to deploy. It's comparable to a marketplace for all tools that can be deployed with an instance of Docktor.
Each template or service contains default information needed by Docktor to automatically deploy a service.

Catalog Service

A catalog service is an entity representing a tool or a service, available for deployment on a given server. A catalog service is versionned and contains one ore more catalog containers. These versions are used to notify the user that a new version is avalaible for his deployed service.
Catalog containers inside a catalog service constitute a functionnaly consistant service. It means that the service is not working properly if one its container is down.

For instance, a service named "SonarQube" is composed of many versions ( 5.6.7, 6.4.1, 6.5, ...). Each version of SonarQube is composed of two catalog containers that are working together: the web application and the database. Each catalog container defines a Docker image version and all default configuration, used to help Docktor to deploy it automatically.

Catalog Container

A catalog container is an entity representing a Docker container configuration. It's not a representation of a deployed Docker container itself, but the representation of the default configuration needed by Docktor to deploy it automatically. It defines mandatory variables, parameters, volumes for the container to work when deployed.

For instance, the container for database of "SonarQube" service defines: the internal 'port' to use, 'variables' like the name of the database, binding 'volumes' for persistent data.

Catalog Template

A catalog template is an entity composed of multiple catalog services. Catalog services in a template are not hard bounded with each other. It's just a set of tools that can deployed together in a just one click. Once, the template is just a boostrap entity meant to deploy multiple services at once. The notion of template disappear once the services are deployed in a group.
That's why templates are not versionned. Templates are refering the latest version of the services that are contained in it.

For instance, a 'software factory' template could contains the catalog services "SonarQube", "Jenkins", "Gitlab" and so on. Once theses services are deployed, they are able to work independently.

Group

A group is a restricted area of isolation where services are deployed. Users can access to a group when they are a member of it or when they are Docktor administrators.
Users can interact with deployed services and containers in a given group (stop, start...). Moderators of a group are able to invite other people to join the group.

Service

A service represents a deployed tool, created from a given catalog service. A service belongs to a given a group. A deployed service represent the actual configuration deployed to a given daemon. Unlike a catalog service, it contains only one version at a time. This version changes when the service is updated to a new one (fetched from catalog).
A service can be stopped, started, restarted, redeployed, deleted  or upgraded.

Container

A container represents a Docker container deployed on a given daemon. It belongs to a service in a group. Unlike a catalog container, it contains actual configuration of web application hosted by the container. A container can be stopped, started, or restarted.
*/
package types
