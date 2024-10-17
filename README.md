# VWAP Web Socket
The goal here is evaluate this code project of VWAP (volume-weighted average price) engine. We are using coinbase websocket feed to stream trade executions and update the VWAP for each trading pair as updates become available.

# Problem Specification:

Retrieve a data feed from the coinbase websocket and subscribe to the matches channel.

- Pull data for the following three trading pairs like:

    - BTC-USD
    - ETH-USD
    - ETH-BTC

- VWAP per trading pair is using a sliding window of 200 data points. Meaning, when a new
  data point arrives through the websocket feed the oldest data point need to fall off and the new one will be
  added such that no more than 200 data points are included in the calculation.
- The first 200 updates will have less than 200 data points included. That’s fine for now!
 
### What is VWAP?
VWAP is a trading benchmark that gives the average price of a security traded throughout some period of time, usually within a day, and it's based on both volume and price. Different from a simple moving average the WVAP price takes into account the level of volume in that security’s trading.
 
### Formula

> LaTex:
> $$\frac{\sum Price*Volume}{\sum Volume}$$ 

> Sum of price multiplied by volume / by the total volume.

## Requirements
- [docker](https://docs.docker.com/get-docker/)
- [make](https://www.gnu.org/software/make/)
 
## Running
```bash
   # on root folder to create a docker image, build and run the app:
   > make start
```
`Take a look on Makefile to other options`

## Output:
- sysout
```
   2022/01/05 06:21:20 Symbol: ETH-USD Trade Sum: 200 VWAP: USD 3813.13
   2022/01/05 06:21:20 Symbol: BTC-USD Trade Sum: 149 VWAP: USD 46312.81
   2022/01/05 06:21:20 Symbol: ETH-USD Trade Sum: 200 VWAP: USD 3813.13
   2022/01/05 06:21:20 Symbol: ETH-USD Trade Sum: 200 VWAP: USD 3813.13
```
 

