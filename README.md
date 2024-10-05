# dns

Set up a sql Server

Setup REST Routes for data in SQL server

Set up Dice DB

According to the TTL of a website, when someone calls it, we need to cache it

Handling of different type -

let's say, if it's A -> send directly
If CName -> recursively query their NS till we find a A record

complete logic in helpers/build_response.go

implement a good logging. Till then we'll use default package only

Seems like all :)
