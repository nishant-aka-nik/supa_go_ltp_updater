https://trendlyne.com/fundamentals/stock-screener/536034/swing-journal-stocks/index/NIFTY500/nifty-500/

- we will add a counter for price if the price crosses the target/stoploss 3 times in 15 min interval it is a confirm target/stoploss hit as the price retained at the level for 45 mins 
- LLD - when the cron will run on the price hit it will update the counter for +1 and an email function will check the count == 3 send the email

entry alerter - 
- when volume is more than 3x than daily average and change percentage more than 3%
- the high of this day will be pivot point
- keep tracking the price against pivot point
- issues - we have to think something about retracement pattern by code 

TODO:
entry filter - 
- volume > 2.5
- change pct > 2.5
- green candle 
- (price_away_from_52high == difference_between_high_and_close) < 2

TODO: 
add all the filtered stocks to a analysis table then we will take the percentage difference between current price and entry price 
we will see the total profit whether it will work or not

TODO: 
- write code to make the "last_traded_price" to be updated using go routines
- write code to upsert data to historic prices table every working day at 5 pm