package main

import (
	"context"
	"crawler-sdk/platforms/douyin"
	"crawler-sdk/platforms/xhs"
	"fmt"
	"log"
	"time"
)

func main() {
	dyTest()
}

func xhsTest() {
	// 平台特定的 cookies
	const xhsCookie = `abRequestId=a9e19a3c-ec62-5d58-82d1-11780180d594; a1=1989d69b653ej0uroayuyvsm4vdy26twup2h6vzan50000384309; webId=936fb47d51c175f4dec0400d43cb1af2; gid=yjYjfKjSDy7jyjYjfKjDKUhu2qdV87Fl0v7vJ0IhM64hfx28vJK9EK888qY4q8j8Y84JiYSf; x-user-id-creator.xiaohongshu.com=60de949b0000000001000c2f; customer-sso-sid=68c51753761449804526765463mdqkwv3gzu3e2d; customerClientId=913136795330800; access-token-creator.xiaohongshu.com=customer.creator.AT-68c5175376144980452676570jcbajjyvvfeitaq; galaxy_creator_session_id=P5r2edI91WGTWafcQqcfekMM05fSN8bK1sAb; galaxy.creator.beaker.session.id=1754987635895010773996; webBuild=4.75.3; web_session=040069b1aad7d3565c79251aa83a4b22b8638f; xsecappid=xhs-pc-web; loadts=1755165211163; acw_tc=0a50885317551652109482088e3f446d3b80b98308dd3884df9ecca40fd6d3; websectiga=3633fe24d49c7dd0eb923edc8205740f10fdb18b25d424d2a2322c6196d2a4ad; sec_poison_id=1003c47d-db6b-4c3d-bc0d-100c0b71ac3c; unread={%22ub%22:%22689d42df000000001b035010%22%2C%22ue%22:%22688aa2bf000000000302629f%22%2C%22uc%22:23}` // 替换为你的小红书 cookie
	// const douyinCookie = "your_douyin_cookie_here"                  // 替换为你的抖音 cookie

	// 示例用户ID
	const xhsUserID = "5d5c36ab0000000001008656"

	// 1. 独立实例化每个平台的客户端
	xhsClient := xhs.New(xhsCookie)
	// douyinClient := douyin.New(douyinCookie)

	// 2. 创建一个带超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 3. 调用小红书的用户信息接口
	fmt.Printf("--- 正在获取小红书用户 [%s] 的信息 ---\n", xhsUserID)
	// userProfile, err := xhsClient.GetUserDetail(ctx, xhsUserID)
	// if err != nil {
	// 	log.Fatalf("获取小红书用户信息失败: %v", err)
	// }

	notes, err := xhsClient.GetUserNotes(ctx, xhsUserID, 0, 10)
	if err != nil {
		log.Fatalf("获取小红书用户笔记失败: %v", err)
	}
	fmt.Println(notes)

}

func dyTest() {
	const dyCookie = `enter_pc_once=1; UIFID_TEMP=8f584679f13361d7674396b25a46f1c75dd7736730eef2bca7f1b3cfea8b216d2fc5e4dff2b610813ba388b3ff59dc5c4af1d20553befcdb23e26214ec33ba5b2b5f72e520048d0beb8cb63fbd5848d2; hevc_supported=true; my_rd=2; bd_ticket_guard_client_web_domain=2; UIFID=8f584679f13361d7674396b25a46f1c75dd7736730eef2bca7f1b3cfea8b216dc61daa8dc016660222357f6c42ce484d57f4531133c856c65db7d65424b7b8e383b5ca1606699759713a78e41218a43a0c5297bb0d450d82ac6fcd0a9a5e5400009aea011b6954469ee7f84b68a096db50250b2e2af1fa2a1f8f907bf94fb3866cd1a77984acd5c3e0ac394c2ddefcf28cd0536450aa851baaf40f4274ddbb34; is_staff_user=false; volume_info=%7B%22volume%22%3A0.97%7D; oc_login_type=LOGIN_PERSON; SelfTabRedDotControl=%5B%5D; passport_mfa_token=Cjd62C2JE4dk5EuBRodFDcYo6e1roN2PZKiTkwnRHgn%2BCNYwMhYknh%2FGugM9D%2BNKdoN32HEd6eHjGkoKPAAAAAAAAAAAAABPS%2FOfnAf%2F4h5UzB7KpogdYZuez0H1rnw99uy9NGTKi1MfFnCWCMxSVLE%2B%2FXAZYG41LRDhj%2FgNGPax0WwgAiIBA2sI27M%3D; d_ticket=a3c045be8b7cbc1d5ee5ab407edecffccaa92; passport_assist_user=CkF6Dw3y2xOTaxor3So4_MarcJmhRpilOSoYy5GZDBhfv-o66ML3BGvMTr4jHS19ciT9EEJZOSQ_JPlDgsUxkLuzohpKCjwAAAAAAAAAAAAAT0u54ZFN2JW7DBZjmBJD89WNEet95NO9fHbrEn-26exFmJ3XZw5YtB_fwHopaPvPh0cQxI_4DRiJr9ZUIAEiAQMdBcf1; n_mh=iWfXbFaDwDEttuTu9ZEBCkODIyB8g1da6hApnhQG96M; sid_guard=d8c0fef7998ca2071758d79ac212bed7%7C1753870469%7C5183999%7CSun%2C+28-Sep-2025+10%3A14%3A28+GMT; uid_tt=9aa2e9814cf019a29548d9f54723a2fc; uid_tt_ss=9aa2e9814cf019a29548d9f54723a2fc; sid_tt=d8c0fef7998ca2071758d79ac212bed7; sessionid=d8c0fef7998ca2071758d79ac212bed7; sessionid_ss=d8c0fef7998ca2071758d79ac212bed7; session_tlb_tag=sttt%7C3%7C2MD-95mMogcXWNeawhK-1__________V1mtlj3sXbkokH0LG5XzdA69gXuWznlsKv3P4vDWySq4%3D; sid_ucp_v1=1.0.0-KDIzM2ZjMDk0NTY1M2UwMzZjNjQ3OTBkNmQ1NDQ0NmQ3MTIyMWQyYmEKIQi4ifD7gfSpAxCF4afEBhjvMSAMMIzyh_AFOAdA9AdIBBoCbGYiIGQ4YzBmZWY3OTk4Y2EyMDcxNzU4ZDc5YWMyMTJiZWQ3; ssid_ucp_v1=1.0.0-KDIzM2ZjMDk0NTY1M2UwMzZjNjQ3OTBkNmQ1NDQ0NmQ3MTIyMWQyYmEKIQi4ifD7gfSpAxCF4afEBhjvMSAMMIzyh_AFOAdA9AdIBBoCbGYiIGQ4YzBmZWY3OTk4Y2EyMDcxNzU4ZDc5YWMyMTJiZWQ3; login_time=1753870470093; _bd_ticket_crypt_cookie=f3715d24cf655fbd73db46fa5c091003; __security_mc_1_s_sdk_sign_data_key_web_protect=b6ba1603-4ea1-aede; __security_server_data_status=1; __security_mc_1_s_sdk_crypt_sdk=978c499b-4ed9-8c07; __security_mc_1_s_sdk_cert_key=5093dedc-470a-938c; passport_csrf_token=db7c040a907b813f87bce98fdcbee44b; passport_csrf_token_default=db7c040a907b813f87bce98fdcbee44b; stream_recommend_feed_params=%22%7B%5C%22cookie_enabled%5C%22%3Atrue%2C%5C%22screen_width%5C%22%3A2048%2C%5C%22screen_height%5C%22%3A1152%2C%5C%22browser_online%5C%22%3Atrue%2C%5C%22cpu_core_num%5C%22%3A12%2C%5C%22device_memory%5C%22%3A8%2C%5C%22downlink%5C%22%3A10%2C%5C%22effective_type%5C%22%3A%5C%224g%5C%22%2C%5C%22round_trip_time%5C%22%3A0%7D%22; FOLLOW_LIVE_POINT_INFO=%22MS4wLjABAAAAb9jM0SSsKo7XkrqdP8jqm_obQoKjhx-5Zzfw_7CAQ3C_NpSzhmjwhJMRLz01R4EC%2F1755360000000%2F0%2F0%2F1755320710342%22; FOLLOW_NUMBER_YELLOW_POINT_INFO=%22MS4wLjABAAAAb9jM0SSsKo7XkrqdP8jqm_obQoKjhx-5Zzfw_7CAQ3C_NpSzhmjwhJMRLz01R4EC%2F1755360000000%2F0%2F1755320110342%2F0%22; home_can_add_dy_2_desktop=%221%22; is_dash_user=1; biz_trace_id=9a430a52; publish_badge_show_info=%221%2C0%2C0%2C1755320135080%22; IsDouyinActive=false; x-web-secsdk-uid=65a4a037-a839-4617-b039-5b98a42a95dd; gfkadpd=2906,33638; _tea_utm_cache_2906=undefined; csrf_session_id=12bf6907a00f6545719fcb949aa9a3d8; bd_ticket_guard_client_data=eyJiZC10aWNrZXQtZ3VhcmQtdmVyc2lvbiI6MiwiYmQtdGlja2V0LWd1YXJkLWl0ZXJhdGlvbi12ZXJzaW9uIjoxLCJiZC10aWNrZXQtZ3VhcmQtcmVlLXB1YmxpYy1rZXkiOiJCUFU1a0E3WkU2Nkp2a0xpM0kwMmt5Rk5vR2NTNjRxTFFkRFAyYXpRUUJGQjRRMlUxRlhJclZTbGJKYU9lSVd2U2tKM3RTOWpQOFFrU1JVWkpzeVBwU289IiwiYmQtdGlja2V0LWd1YXJkLXdlYi12ZXJzaW9uIjoyfQ%3D%3D; ttwid=1%7CGb-ju4M-tlnkctdysXeQBdDagZ_KOyU-47QhJ3n3lxo%7C1755322443%7Cd18434f9d37dbb75d487f3d71d6db5b72d679e9cebcc454974660a4172565fbe; odin_tt=cf1d52a35359648e83804d98bffac398408cf78c49049ee05bf3656f5990bcf9a6f822befb49094819ac418f9b722a12ce58a0c7ba4b36914a6ffcdfa8a5e828; _tea_utm_cache_1128=undefined; passport_fe_beating_status=true; gd_random=eyJtYXRjaCI6dHJ1ZSwicGVyY2VudCI6MC44MDYxNDQ3MjE0MDk0MDczfQ==.cNqkKv1qJ/8BVjZMVPr3Sy1434h3NHusfnQUD6fo3tg=`
	dyClient := douyin.New(dyCookie)
	// 2. 创建一个带超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	dyClient.PublishImage(ctx, "test.jpg", "test", "test")

}
