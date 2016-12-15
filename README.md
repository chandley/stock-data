Shows simple price data for US stocks traded on the main exchange

It's running on heroku at https://go-stock-data.herokuapp.com/composite/F
use the endpoint with the ticker for the stock you want to see (eg FB, MSFT)

to get data it needs API keys for Quandl (price data) and Edgar (fundamentals) 
https://www.quandl.com/
http://developer.edgar-online.com/

these should be set as environmental variables

QUANDL_API_KEY
EDGAR_API_KEY