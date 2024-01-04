package aicontroller

import (
	"fmt"
	"net/http"

	"github.com/ayan-sh03/triviagenious-backend/internal/util"
)

// 1. take file input
// 2. convert it into text
/// 3. feed it into the api

func GetQuestionFromAi(w http.ResponseWriter, r *http.Request) {

	// text := `The stock market is a dynamic and intricate financial ecosystem where investors buy and sell shares of publicly traded companies. It serves as a platform for companies to raise capital by issuing stocks, and investors, in turn, gain ownership in those companies. Stock prices fluctuate based on various factors, including economic indicators, company performance, and market sentiment. Investors often analyze financial data, market trends, and company fundamentals to make informed decisions, aiming to capitalize on potential gains or mitigate losses. The stock market plays a crucial role in the broader economy, influencing consumer confidence, corporate behavior, and overall economic stability.

	//  Investing in the stock market involves inherent risks and uncertainties. Stock prices are subject to market volatility, influenced by global events, economic indicators, and geopolitical developments. Investors employ various strategies, such as value investing, technical analysis, and diversification, to navigate market fluctuations and optimize their investment portfolios. Market participants include individual investors, institutional investors, and traders who engage in buying and selling activities on exchanges worldwide. The stock market reflects the collective perceptions and expectations of investors, contributing to the continuous evolution of financial markets.

	//  The stock market is a key barometer of economic health, serving as a mirror that reflects the underlying strength and performance of businesses and industries. It provides a mechanism for capital allocation, enabling companies to grow and innovate. Additionally, stock market indices, such as the Dow Jones Industrial Average and the S&P 500, are widely used as benchmarks to assess the overall health of the financial markets and gauge economic trends. Overall, the stock market is a complex and integral component of the global financial landscape, influencing wealth creation, economic development, and investment strategies worldwide.`

	response := util.CopyFileFromRequest(w, r)

	if response == nil {
		util.RespondWithError(w, http.StatusInternalServerError, "Error Reading File")
		return
	}

	extractedText := util.ExtractTextFromPdf(response.Name())

	fmt.Println(extractedText)
	res := util.ExecuteQuery(extractedText)

	if res == "" {
		util.RespondWithError(w, http.StatusInternalServerError, "Could not fetch from AI")
		return
	}

	util.RespondWithJSON(w, http.StatusOK, res)

}
