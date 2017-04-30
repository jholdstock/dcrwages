require 'net/http'
require 'json'
require 'date' 
require 'logger' 

class Array;
	def sum; inject(nil) { |sum,x| sum ? sum + x : x }; end;
	def avg; sum / size; end; 
end

$l = Logger.new(STDOUT)
$l.level = Logger::INFO
$l.formatter = proc do |severity, date, p, msg|
   "[#{severity}] #{msg}\n"
   #"#{date.strftime('%Y-%m-%d %H:%M:%S')}: [#{severity}] #{msg}\n"
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

	$l.debug(uri.request_uri)

	http = Net::HTTP.new(uri.host, uri.port)
	http.use_ssl = true
	request = Net::HTTP::Get.new(uri.request_uri)
	response = http.request(request)

	sleep 0.2 # Polo doesn't like >6 requests per second

	# TODO: Trusting polo to return sorted data. Should sort here to be safe
	JSON.parse(response.body)
end	

def get_month_avg mon, year
	start_date = Date.new(year, mon, 1)
	end_date = start_date.next_month

	btc = download_chart(start_date, end_date, "BTC_DCR")
	usd = download_chart(start_date, end_date, "USDT_BTC")

	dcr_usd = []

	[btc.length, usd.length].min.times do |i|
		btc_avg = [btc[i]["high"], btc[i]["low"]].avg
		usd_avg = [usd[i]["high"], usd[i]["low"]].avg

		dcr_usd.push(btc_avg * usd_avg)
	end

	month_avg = sprintf('%.3f', dcr_usd.avg)
	$l.info "#{month_avg} USD/DCR (#{mon}-#{year})"
end

get_month_avg(12, 2016)
get_month_avg(1, 2017)
get_month_avg(2, 2017)
get_month_avg(3, 2017)
get_month_avg(4, 2017)
