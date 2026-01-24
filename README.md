# url-shortener

This is a url shortener which generates shorten urls and we are also storing some analytics (mainly how many clicks this shorturl has received)

**logic** : generating an 8 digit unique id for url (0-9 , a-z and A-Z) so 62 characters , 62^8 gives which is 218 trillions so a lot of urls.
**main drawback**  : same urls can still generate different unique ids I have not yet taken into account this logic yet so we might have entries in db like.

google.com -> asdj123n

google.com -> xAnasd12

and so on.

A db check for this could be costly (could add an index on original url) but did not thought about this that much tbh.

Also in the db schema we have a timestamp at which we have created the url, curretly this is useless.

In future I'm thinking of adding the feature of link expired.

Might use Hadoop for this to perform a batch job on a per day basis.

**Mainly has 2 endpoints**

**POST** localhost:4000/shorten

{

"url" : "http://google.com"

}

gives back response

{"short_url":"http://localhost:4000/r/OvXMU4cf"}

**GET** localhost:4000/r/{id} these url that we got from above

gives back the orignal url

{"original_url":"http://www.google.com","source":"cache/database"}

Have used Redis caching in this for get call.

For the Analytics part we are just storing the clicks ,but since as per the logic everytime we are doing a get call we are basically writing in the database

to increment the counter of clicks which would mean same amount of reads and writes.

The Analytics need not to be consistent so they are updated Asynchronously.

All the events of clicks are added in Kafka and then we will use Apache Spark to perform mini batch Aggregation of clicks and then update in the database after every minute interval (Spark Consumer is not yet been implemented so curretly in db you will see only 1 click per url)



