package playerinfo

import (
	"net"

	. "github.com/KouKouChan/CSO2-Server/blademaster/typestruct"
	. "github.com/KouKouChan/CSO2-Server/kerlong"
	. "github.com/KouKouChan/CSO2-Server/servermanager"
	. "github.com/KouKouChan/CSO2-Server/verbose"
)

const (
	started  = 0
	finished = 1
)

func OnSetCampaign(p *PacketData, client net.Conn) {
	var pkt InSetCampaignPacket
	if !p.PraseSetCampaignPacket(&pkt) {
		DebugInfo(2, "Error : Client from", client.RemoteAddr().String(), "sent a illegal SetCampaign packet !")
		return
	}
	//找到对应用户
	uPtr := GetUserFromConnection(client)
	if uPtr == nil ||
		uPtr.Userid <= 0 {
		DebugInfo(2, "Error : Client from", client.RemoteAddr().String(), "try to SetCampaign but not in server !")
		return
	}
	//修改数据
	switch pkt.PacketType {
	case started:
		DebugInfo(1, "User", string(uPtr.UserName), "Started Campaign")
	case finished:
		if isMissionCampaignIdValid(pkt.CampaignId) {
			//发送数据包
			rst := BytesCombine(BuildHeader(uPtr.CurrentSequence, PacketTypeUserInfo), BuildUserInfo(0x1000, NewUserInfo(uPtr), uPtr.Userid, true))
			SendPacket(rst, uPtr.CurrentConnection)
			DebugInfo(1, "User", string(uPtr.UserName), "Updated Campaign ", pkt.CampaignId)
		} else {
			DebugInfo(2, "User", string(uPtr.UserName), "sent a invalid SetCampaign packet id", pkt.CampaignId)
		}
	default:
		DebugInfo(2, "User", string(uPtr.UserName), "sent a unkown SetCampaign packet type", pkt.PacketType)
	}
}

func isMissionCampaignIdValid(id uint16) bool {
	return id == Campaign_0 ||
		id == Campaign_1 ||
		id == Campaign_2 ||
		id == Campaign_3 ||
		id == Campaign_4 ||
		id == Campaign_5 ||
		id == Campaign_6
}
