package ucs

import (
    "strings"
)

var ucsBodyData map[string]string

func init() {
    ucsBodyData = make(map[string]string)
    fillUCSBodyData()
}

func getAllMacPools(url string, cookie string) string {
	return getData(url, xmlReplaceString(ucsBodyData["allMACPools"],"COOKIE",cookie))
}

func getAllServers(url string, cookie string) string {
    return getData(url, xmlReplaceString(ucsBodyData["allServers"],"COOKIE",cookie))
}

func getUnassociatedServers(url string, cookie string) string {
    return getData(url, xmlReplaceString(ucsBodyData["unassociatedServers"],"COOKIE",cookie))
}

func getAssociatedServers(url string, cookie string) string {
    return getData(url, xmlReplaceString(ucsBodyData["associatedServers"],"COOKIE",cookie))
}

func getAllFaults(url string, cookie string) string {
    return getData(url, xmlReplaceString(ucsBodyData["allFaults"], "COOKIE",cookie))
}

func getAllBlades(url string, cookie string) string {
    return getData(url, xmlReplaceString(ucsBodyData["allServers"], "COOKIE", cookie))
}

func getCPUDetail(url string, cookie string, xml string) string {
    return getData(url, xmlReplaceStringArray(ucsBodyData["serverCPUDetail"],[]string{"COOKIE","XML"},[]string{cookie,xml}))
}

func getDemoData(body string) string {
    if strings.Contains(body, "aaaLogin") {
        return "<aaaLogin cookie=\"\" response=\"yes\" outCookie=\"1476996526/478783f1-01cd-4de3-a3c0-d78a118f842f\" outRefreshPeriod=\"600\" outPriv=\"admin,read-only\" outDomains=\"\" outChannel=\"noencssl\" outEvtChannel=\"noencssl\" outSessionId=\"web_60426_B\" outVersion=\"2.2(3e)\" outName=\"admin\"> </aaaLogin>"
    } else if strings.Contains(body, "lsServer") {
        tmp := "<configResolveClass cookie=\"1476996526/478783f1-01cd-4de3-a3c0-d78a118f842f\" response=\"yes\" classId=\"lsServer\"><outConfigs>"
        tmp += "<lsServer agentPolicyName=\"\" assignState=\"assigned\" assocState=\"associated\" biosProfileName=\"PLAY-BIOS\" bootPolicyName=\"Windows-Boot\" configQualifier=\"\" configState=\"applied\" descr=\"\" dn=\"org-root/org-Core/org-Windows/ls-WIN-BM-SP-1\" dynamicConPolicyName=\"\" extIPPoolName=\"ext-mgmt\" extIPState=\"none\" fltAggr=\"0\" fsmDescr=\"\" fsmFlags=\"\" fsmPrev=\"ConfigureSuccess\" fsmProgr=\"100\" fsmRmtInvErrCode=\"none\" fsmRmtInvErrDescr=\"\" fsmRmtInvRslt=\"\" fsmStageDescr=\"\" fsmStamp=\"2016-01-25T14:25:14.701\" fsmStatus=\"nop\" fsmTry=\"0\" hostFwPolicyName=\"HFP-2.2.1b\" identPoolName=\"UKDCV-PLAY-UUID\" intId=\"5354557\" kvmMgmtPolicyName=\"\" localDiskPolicyName=\"Windows-LDCP\" maintPolicyName=\"PLAY-MAINT\" mgmtAccessPolicyName=\"\" mgmtFwPolicyName=\"\" name=\"WIN-BM-SP-1\" operBiosProfileName=\"org-root/bios-prof-PLAY-BIOS\" operBootPolicyName=\"org-root/org-Core/org-Windows/boot-policy-Windows-Boot\" operDynamicConPolicyName=\"\" operExtIPPoolName=\"\" operHostFwPolicyName=\"org-root/fw-host-pack-HFP-2.2.1b\" operIdentPoolName=\"org-root/uuid-pool-UKDCV-PLAY-UUID\" operKvmMgmtPolicyName=\"\" operLocalDiskPolicyName=\"org-root/org-Core/org-Windows/local-disk-config-Windows-LDCP\" operMaintPolicyName=\"org-root/maint-PLAY-MAINT\" operMgmtAccessPolicyName=\"\" operMgmtFwPolicyName=\"\" operPowerPolicyName=\"org-root/power-policy-PLAY-PCP\" operScrubPolicyName=\"org-root/scrub-default\" operSolPolicyName=\"\" operSrcTemplName=\"org-root/org-Core/org-Windows/ls-Windows-BM-TEMP\" operState=\"ok\" operStatsPolicyName=\"org-root/thr-policy-default\" operVconProfileName=\"\" operVmediaPolicyName=\"\" owner=\"management\" pnDn=\"sys/chassis-4/blade-4\" policyLevel=\"0\" policyOwner=\"local\" powerPolicyName=\"PLAY-PCP\" resolveRemote=\"yes\" scrubPolicyName=\"\" solPolicyName=\"\" srcTemplName=\"Windows-BM-TEMP\" statsPolicyName=\"default\" svnicConfig=\"yes\" type=\"instance\" usrLbl=\"\" uuid=\"4b279c8a-5f19-11e2-0000-0000000003bf\" uuidSuffix=\"0000-0000000003BF\" vconProfileName=\"\" vmediaPolicyName=\"\"/>"
        tmp += "</outConfigs></configResolveClass>"
        return tmp
    }
    return ""
}

func login(url string, username string, password string) string {
    return getData(url, xmlReplaceStringArray(ucsBodyData["login"],[]string{"USERNAME","PASSWORD"},[]string{username,password}))
}

func logout(url string, cookie string) string {
    return getData(url, xmlReplaceString(ucsBodyData["logout"],"COOKIE",cookie))
}

func fillUCSBodyData() {
    ucsBodyData["login"]                            = "<aaaLogin inName=\"#USERNAME#\" inPassword=\"#PASSWORD#\"/>"
    ucsBodyData["logout"]                           = "<aaaLogout inCookie=\"#COOKIE#\" />"
    ucsBodyData["associatedServers"]                = "<configResolveClass cookie=\"#COOKIE#\" inHierarchical=\"false\" classId=\"lsServer\"><inFilter><eq class=\"lsServer\" property=\"assocState\" value=\"associated\" /></inFilter></configResolveClass>"
    ucsBodyData["unassociatedServers"]              = "<configResolveClass cookie=\"#COOKIE#\" inHierarchical=\"false\" classId=\"lsServer\"><inFilter><eq class=\"lsServer\" property=\"assocState\" value=\"unassociated\" /></inFilter></configResolveClass>"
    ucsBodyData["allServiceProfiles"]               = "<configResolveClass cookie=\"#COOKIE#\" inHierarchical=\"false\" classId=\"lsServer\"></configResolveClass>"
    ucsBodyData["allMACPools"]                      = "<configScope cookie=\"#COOKIE#\" inHierarchical=\"false\" dn=\"mac\" inClass=\"macpoolAddr\"/>"
    ucsBodyData["allFaults"]                        = "<configResolveClass cookie=\"#COOKIE#\" inHierarchical=\"false\" classId=\"faultInst\"/>"
    ucsBodyData["allServersWithGreaterMemory"]      = "<configResolveClass cookie=\"#COOKIE#\" inHierarchical=\"false\" classId=\"memoryArray\"><inFilter><gt class=\"memoryArray\" property=\"currCapacity\" value=\"#VALUE#\" /></inFilter></configResolveClass>"
    ucsBodyData["allServersWithLessMemory"]         = "<configResolveClass cookie=\"#COOKIE#\" inHierarchical=\"false\" classId=\"memoryArray\"><inFilter><lt class=\"memoryArray\" property=\"currCapacity\" value=\"#VALUE#\"/></inFilter></configResolveClass>"
    ucsBodyData["serverDetail"]                     = "<configResolveDn cookie=\"#COOKIE#\" dn=\"#DN#\" inHierarchical=\"false\"/>"
    ucsBodyData["chassisDetail"]                    = "<configResolveClass cookie=\"#COOKIE#\" classId=\"equipmentChassis\" inHierarchical=\"false\"><inFilter></inFilter></configResolveClass>"
    ucsBodyData["allServers"]                       = "<configResolveClass cookie=\"#COOKIE#\" classId=\"computeBlade\" inHierarchical=\"false\"><inFilter></inFilter></configResolveClass>"
    ucsBodyData["serverDeepCPUDetail"]              = "<configScope cookie=\"#COOKIE#\" inHierarchical=\"false\" dn=\"#NODE#/board/#CPU#\" inClass=\"processorEnvStats\"/>"
    ucsBodyData["findAllType"]                      = "<configFindDnsByClassId cookie=\"#COOKIE#\" classId=\"#TYPE#\" />"
    ucsBodyData["getFIAttributes"]                  = "<configResolveClass cookie=\"#COOKIE#\" classId=\"networkElement\" inHierarchical=\"false\"><inFilter></inFilter></configResolveClass>"
    ucsBodyData["getFIFirmware"]                    = "<configResolveChildren cookie=\"#COOKIE#\" classId=\"firmwareBootUnit\" inDn='sys/switch-#FI#/mgmt/fw-boot-def\" inHierarchical=\"false\"><inFilter></inFilter></configResolveChildren>"
    ucsBodyData["getFIDeviceInfo"]                  = "<configResolveClass cookie=\"#COOKIE#\" classId=\"equipmentSwitchCard\" inHierarchical=\"false\"><inFilter></inFilter></configResolveClass>"
    ucsBodyData["getFIPortInfo"]                    = "<configResolveClass cookie=\"#COOKIE#\" classId=\"etherPIo\" inHierarchical=\"false\"><inFilter></inFilter></configResolveClass>"
    ucsBodyData["getFIPSUInfo"]                     = "<configResolveClass cookie=\"#COOKIE#\" classId=\"equipmentPsu\" inHierarchical=\"false\"><inFilter></inFilter></configResolveClass>"
    ucsBodyData["getFIDeviceStats"]                 = "<configResolveChildren cookie=\"#COOKIE#\" classId=\"swSystemStats\" inDn=\"sys/switch-#FI#\" inHierarchical=\"false\"><inFilter></inFilter></configResolveChildren>"
    ucsBodyData["getFIStats"]                       = "<configScope cookie=\"#COOKIE#\" inClass=\"swSystemStatsHist\" inHierarchical=\"false\" dn=\"sys/switch-#FI#\" ></configScope>"
    ucsBodyData["serverCPUDetail"]                  = "<configResolveDns cookie=\"#COOKIE#\" inHierarchical=\"false\"><inDns>#XML#</inDns></configResolveDns>"
}
