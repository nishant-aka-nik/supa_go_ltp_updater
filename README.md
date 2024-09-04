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

TODO:
3 stage filtering

- Alert
    - has to be green candle
    - query stock which are active
    - if price is in between 2% of cross_match_pivot and 5% of cross_match_pivot and active == true
    - active means it is near about equal to 52 week high
    - before sending the alert mark the flag entry as true
    - if not in entry band mark it as false for the UI top picks page 
    - top picks will only show entry stocks
    - volume > 1.2
- Insert
    - if cross_match == true of todays data insert the record with active = true and cross_match_pivot = ltp 
    - and also check that this record is not present in the table and active != true
    - meaning before inserting query all the records which are active and skip insertion
    - volume > 1.5
- Reset 
    - if the active == true and red candle then mark the active = false

- filter cron to start everyday at 2

- it doesnt have to be accurate match cause we might miss good ones

Changelog 24 Aug 2024: New Ideas
 - More Dynamic Stoploss
 - top picks should have 2 columns 
    - active which should be not more than 3 days like 3 din tk to entryable rahega 
    - uske baad vo chala jaega old top picks me
 - health check mails

DB Design
- https://dbdiagram.io/d/Argus-6664a4a89713410b051854cf


BUG TODO:
- reset entry based on date also 
- health check report needs to be generated and then sent properly only once per day
- also test the reset of the stock
    - this fucntion will operate on filter_history table
    - check that what will happen if the stocks becomes active again will the new entry be create or same entry gets updated
- new target implementation
- add cross matched stocks list with progress bar to see how far it is from entry price
    - 2 sections will be there namely - potential entry stocks and Ready to entry stocks
- add a new table 'backtester' which will add stocks and use my 5target-10stoploss-15stoploss strategy to mark entry and exit 
    - entry bool
    - entry_date date
    - exit when stoploss hit 
    - exit bool
    - exit_date date
    - what will happen if exited stock becomes active again
    - i will add a new fucntion named backtester which will run in the end of the stock market session daily
    - I am being lenient on the accurate stoploss lets see how much i can get by this
    - this new fucntion backtester will operate on new table backtester
- we are decoupling back test results and top picks 
