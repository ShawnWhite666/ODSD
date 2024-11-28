package interfaces

import (
	"bytes"
	"encoding/json"
	"github.com/BitofferHub/lotterysvr/internal/biz"
	"github.com/BitofferHub/lotterysvr/internal/constant"
	"io"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestAddPrizeList(t *testing.T) {
	client := &http.Client{}
	addPrize1 := biz.ViewPrize{
		Title:     "iphone",
		Img:       "https://p0.ssl.qhmsg.com/t016ff98b934914aca6.png",
		PrizeNum:  10,
		PrizeCode: "1-10",
		EndTime:   time.Now().Add(time.Hour * 24 * 7),
		BeginTime: time.Now(),
		PrizeType: constant.PrizeTypeEntityLarge,
	}

	addPrize2 := biz.ViewPrize{
		Title:     "homepod",
		Img:       "https://imgservice.suning.cn/uimg1/b2c/image/t_QerWgoH9ergm0_NY4WhA.png_800w_800h_4e",
		PrizeNum:  50,
		PrizeCode: "100-150",
		EndTime:   time.Now().Add(time.Hour * 24 * 7),
		BeginTime: time.Now(),
		PrizeType: constant.PrizeTypeEntityMiddle,
	}

	addPrize3 := biz.ViewPrize{
		Title:     "充电器",
		Img:       "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcS5Y7iNXpcthgJ6yvE3Os1bTwLARVnvwYXeKA&usqp=CAU",
		PrizeNum:  500,
		PrizeCode: "1500-2000",
		EndTime:   time.Now().Add(time.Hour * 24 * 7),
		BeginTime: time.Now(),
		PrizeType: constant.PrizeTypeEntitySmall,
	}

	addPrize4 := biz.ViewPrize{
		Title:     "优惠券",
		Img:       "https://static.699pic.com/images/diversion/d66d647c52cd66beb800ba09748ea080.jpgU",
		PrizeNum:  8000,
		PrizeCode: "2001-9999",
		EndTime:   time.Now().Add(time.Hour * 24 * 7),
		BeginTime: time.Now(),
		PrizeType: constant.PrizeTypeCouponDiff,
	}
	prizeList := []*biz.ViewPrize{
		&addPrize1, &addPrize2, &addPrize3, &addPrize4,
	}

	addPrizeListReq := AddPrizeListReq{
		UserID:    1,
		PrizeList: prizeList,
	}

	bytesData, err := json.Marshal(&addPrizeListReq)
	if err != nil {
		t.Errorf("Error marshalling:%v", err)
	}
	t.Logf("req json = %s\n", string(bytesData))
	req, _ := http.NewRequest("POST", "http://localhost:10080/admin/add_prize_list", bytes.NewReader(bytesData))
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	body, _ := io.ReadAll(resp.Body)
	bodystr := string(body)
	if err != nil {
		t.Errorf("add prize list http request err:%v\n", err)
	}
	t.Logf("rspStr=%s\n", bodystr)

}

func TestClearPrize(t *testing.T) {
	clearPrizeReq := ClearPrizeReq{
		UserID: 1,
	}
	bytesData, err := json.Marshal(&clearPrizeReq)
	if err != nil {
		t.Errorf("Error marshalling:%v", err)
	}
	t.Logf("req json = %s\n", string(bytesData))
	req, _ := http.NewRequest("POST", "http://localhost:10080/admin/clear_prize", bytes.NewReader(bytesData))
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	body, _ := io.ReadAll(resp.Body)
	bodystr := string(body)
	if err != nil {
		t.Errorf("clear prize http request err:%v\n", err)
	}
	t.Logf("rspStr=%s\n", bodystr)
}

func TestImportCoupon(t *testing.T) {
	codeStr := ""
	for i := 1; i <= 50; i++ {
		code := "coupon_code00000"
		if i < 10 {
			code = code + "0" + strconv.Itoa(i) + "\n"
		} else {
			code = code + strconv.Itoa(i) + "\n"
		}
		codeStr = codeStr + code
	}
	couponCode := strings.Trim(codeStr, "\n")
	couponInfo := biz.ViewCouponInfo{
		PrizeId:    4,
		Code:       couponCode,
		SysCreated: time.Time{},
		SysUpdated: time.Time{},
		SysStatus:  1, // 正常
	}
	importCouponReq := ImportCouponReq{
		UserID:     1,
		CouponInfo: &couponInfo,
	}
	bytesData, err := json.Marshal(&importCouponReq)
	if err != nil {
		t.Errorf("Error marshalling:%v", err)
	}
	t.Logf("req json = %s\n", string(bytesData))
	req, _ := http.NewRequest("POST", "http://localhost:10080/admin/import_coupon", bytes.NewReader(bytesData))
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	body, _ := io.ReadAll(resp.Body)
	bodystr := string(body)
	if err != nil {
		t.Errorf("import coupon http request err:%v\n", err)
	}
	t.Logf("rspStr=%s\n", bodystr)
}

func TestImportCouponWithCache(t *testing.T) {
	codeStr := ""
	for i := 1; i <= 50; i++ {
		code := "coupon_code00000"
		if i < 10 {
			code = code + "0" + strconv.Itoa(i) + "\n"
		} else {
			code = code + strconv.Itoa(i) + "\n"
		}
		codeStr = codeStr + code
	}
	couponCode := strings.Trim(codeStr, "\n")
	couponInfo := biz.ViewCouponInfo{
		PrizeId:    73,
		Code:       couponCode,
		SysCreated: time.Time{},
		SysUpdated: time.Time{},
		SysStatus:  1, // 正常
	}
	importCouponReq := ImportCouponReq{
		UserID:     1,
		CouponInfo: &couponInfo,
	}
	bytesData, err := json.Marshal(&importCouponReq)
	if err != nil {
		t.Errorf("Error marshalling:%v", err)
	}
	t.Logf("req json = %s\n", string(bytesData))
	req, _ := http.NewRequest("POST", "http://localhost:10080/admin/import_coupon_cache", bytes.NewReader(bytesData))
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	body, _ := io.ReadAll(resp.Body)
	bodystr := string(body)
	if err != nil {
		t.Errorf("import coupon with cache http request err:%v\n", err)
	}
	t.Logf("rspStr=%s\n", bodystr)
}

func TestClearCoupon(t *testing.T) {
	clearCouponReq := ClearCouponReq{
		UserID: 1,
	}
	bytesData, err := json.Marshal(&clearCouponReq)
	if err != nil {
		t.Errorf("Error marshalling:%v", err)
	}
	t.Logf("req json = %s\n", string(bytesData))
	req, _ := http.NewRequest("POST", "http://localhost:10080/admin/clear_coupon", bytes.NewReader(bytesData))
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	body, _ := io.ReadAll(resp.Body)
	bodystr := string(body)
	if err != nil {
		t.Errorf("clear coupon http request err:%v\n", err)
	}
	t.Logf("rspStr=%s\n", bodystr)
}

func TestClearLotteryTimes(t *testing.T) {
	clearLotteryTimesReq := ClearLotteryTimesReq{
		UserID: 1,
	}
	bytesData, err := json.Marshal(&clearLotteryTimesReq)
	if err != nil {
		t.Errorf("Error marshalling:%v", err)
	}
	t.Logf("req json = %s\n", string(bytesData))
	req, _ := http.NewRequest("POST", "http://localhost:10080/admin/clear_lottery_times", bytes.NewReader(bytesData))
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	body, _ := io.ReadAll(resp.Body)
	bodystr := string(body)
	if err != nil {
		t.Errorf("clear lottery times http request err:%v\n", err)
	}
	t.Logf("rspStr=%s\n", bodystr)
}

func TestClearResult(t *testing.T) {
	clearResultReq := ClearResultReq{
		UserID: 1,
	}
	bytesData, err := json.Marshal(&clearResultReq)
	if err != nil {
		t.Errorf("Error marshalling:%v", err)
	}
	t.Logf("req json = %s\n", string(bytesData))
	req, _ := http.NewRequest("POST", "http://localhost:10080/admin/clear_result", bytes.NewReader(bytesData))
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	body, _ := io.ReadAll(resp.Body)
	bodystr := string(body)
	if err != nil {
		t.Errorf("clear result http request err:%v\n", err)
	}
	t.Logf("rspStr=%s\n", bodystr)
}

func TestLotteryV1(t *testing.T) {
	lotteryReq := LotteryReq{
		UserID:   1,
		UserName: "zhangsan",
		IP:       "192.168.9.9",
	}
	bytesData, err := json.Marshal(&lotteryReq)
	if err != nil {
		t.Errorf("Error marshalling:%v", err)
	}
	req, _ := http.NewRequest("POST", "http://localhost:10080/lottery/v1/get_lucky", bytes.NewReader(bytesData))
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	body, _ := io.ReadAll(resp.Body)
	bodystr := string(body)
	if err != nil {
		t.Errorf("lotteryv1 http request err:%v\n", err)
	}
	t.Logf("rspStr=%s\n", bodystr)
}

func TestLotteryV2(t *testing.T) {
	lotteryReq := LotteryReq{
		UserID:   2,
		UserName: "lisi",
		IP:       "192.168.9.10",
	}
	bytesData, err := json.Marshal(&lotteryReq)
	if err != nil {
		t.Errorf("Error marshalling:%v", err)
	}
	req, _ := http.NewRequest("POST", "http://localhost:10080/lottery/v2/get_lucky", bytes.NewReader(bytesData))
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	body, _ := io.ReadAll(resp.Body)
	bodystr := string(body)
	if err != nil {
		t.Errorf("lotteryv2 http request err:%v\n", err)
	}
	t.Logf("rspStr=%s\n", bodystr)
}

func TestLotteryV3(t *testing.T) {
	lotteryReq := LotteryReq{
		UserID:   2,
		UserName: "lisi",
		IP:       "192.168.9.10",
	}
	bytesData, err := json.Marshal(&lotteryReq)
	if err != nil {
		t.Errorf("Error marshalling:%v", err)
	}
	req, _ := http.NewRequest("POST", "http://localhost:10080/lottery/v3/get_lucky", bytes.NewReader(bytesData))
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	body, _ := io.ReadAll(resp.Body)
	bodystr := string(body)
	if err != nil {
		t.Errorf("lotteryv3 http request err:%v\n", err)
	}
	t.Logf("rspStr=%s\n", bodystr)
}
