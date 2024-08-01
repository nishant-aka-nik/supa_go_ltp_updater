https://trendlyne.com/fundamentals/stock-screener/536034/swing-journal-stocks/index/NIFTY500/nifty-500/

- we will add a counter for price if the price crosses the target/stoploss 3 times in 15 min interval it is a confirm target/stoploss hit as the price retained at the level for 45 mins 
- LLD - when the cron will run on the price hit it will update the counter for +1 and an email function will check the count == 3 send the email

entry alerter - 
- when volume is more than 3x than daily average and change percentage more than 3%
- the high of this day will be pivot point
- keep tracking the price against pivot point
- issues - we have to think something about retracement pattern by code 

TODO: Done
entry filter - 
- volume > 1.5
- cross match
- openCloseDiff > 2
- highCloseDiff < 2.5

TODO: 
 - add all the filtered stocks to a analysis table then we will take the percentage difference between current price and entry price 
 - we will analyse the total profit whether it will work or not
 - add new column - profit, in filter_history and argus will calculate the profit every day on it 
 - add new column trade_completed bool it will also mark the exit of the stock when price cross down the 20 ema and we will not calculate the profit on these stocks

TODO: 
- write code to make the "last_traded_price" to be updated using go routines
- write code to upsert data to historic prices table every working day at 5 pm
- write a code to calculate 20 ema of the stock 