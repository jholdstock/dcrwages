require 'net/http'
require 'json'
require 'date'

class Array;
	def sum; inject(nil) { |sum,x| sum ? sum + x : x }; end;
	def avg; sum / size; end; 
end

def download_chart start_date, end_date, currency_pair
	url = "https://poloniex.com/public"
	uri = URI.parse(url)

	uri.query = URI.encode_www_form({ 
		"command"      => "returnChartData",
		"currencyPair" => currency_pair,
		"start"        => start_date.to_time.to_i,
		"end"          => end_date.to_time.to_i,
		"period"       => 900,
	})

	# puts uri.request_uri
	
	http = Net::HTTP.new(uri.host, uri.port)
	http.use_ssl = true
	request = Net::HTTP::Get.new(uri.request_uri)
	response = http.request(request)

	
	# TODO: Trusting polo to return sorted data. Should sort here to be safe
	JSON.parse(response.body)
end	

def get_month_avg mon, year
	start_date = DateTime.new(year, mon, 1)
	end_date = start_date.next_month
	
	btc = download_chart(start_date, end_date, "BTC_DCR")
	usdt = download_chart(start_date, end_date, "USDT_BTC")
	
	dcr_usdt = []
	
	[btc.length, usdt.length].min.times do |i|
		btc_avg = btc[i]["weightedAverage"]
		usdt_avg = usdt[i]["weightedAverage"]
		
		dcr_usdt.push(btc_avg * usdt_avg)
	end
	
	month_avg = sprintf('%.4f', dcr_usdt.avg)
	puts "#{month_avg} USD/DCR (#{mon}-#{year})"
	sleep 0.3 # Polo doesn't like >6 requests per second
end

get_month_avg(6, 2016)
get_month_avg(7, 2016)
get_month_avg(8, 2016)
get_month_avg(9, 2016)
get_month_avg(10, 2016)
get_month_avg(11, 2016)
get_month_avg(12, 2016)

get_month_avg(1, 2017)
get_month_avg(2, 2017)
get_month_avg(3, 2017)
get_month_avg(4, 2017)
get_month_avg(5, 2017)
get_month_avg(6, 2017)
get_month_avg(7, 2017)
get_month_avg(8, 2017)
get_month_avg(9, 2017)
get_month_avg(10, 2017)
get_month_avg(11, 2017)
get_month_avg(12, 2017)

get_month_avg(1, 2018)
get_month_avg(2, 2018)
get_month_avg(3, 2018)
get_month_avg(4, 2018)
get_month_avg(5, 2018)
get_month_avg(6, 2018)
get_month_avg(7, 2018)
get_month_avg(8, 2018)
get_month_avg(9, 2018)
get_month_avg(10, 2018)
get_month_avg(11, 2018)
get_month_avg(12, 2018)

get_month_avg(1, 2019)
get_month_avg(2, 2019)
