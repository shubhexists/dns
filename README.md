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

start docker by
`docker compose --env-file .env up -d`

---

Things which are intentionally not implemented

Obselete QTypes like
QTYPE_MD = 3 // MD (Mail Destination, obsolete)
QTYPE_MF = 4 // MF (Mail Forwarder, obsolete)
QTYPE_MB = 7 // MB (Mailbox Domain Name, obsolete)
QTYPE_MG = 8 // MG (Mail Group Member, obsolete)
QTYPE_MR = 9 // MR (Mail Rename Domain Name, obsolete)

coz why care handling if they aren't used :D

---

also, Class is just Internet class. No other classes of now


changes to make in the last - 
combine all the distributed Flags into one FLag in the Header struct