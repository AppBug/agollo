package agollo

import (
	"testing"
	"github.com/zouyx/agollo/test"
	"time"
	"fmt"
)

//func TestInitRefreshInterval(t *testing.T) {
//	refresh_interval=1*time.Second
//
//	var c AbsComponent
//	c=&AutoRefreshConfigComponent{}
//	c.Start()
//}

func TestAutoSyncConfigServices(t *testing.T) {
	go runMockConfigServer(normalConfigResponse)
	defer closeMockConfigServer()

	time.Sleep(1*time.Second)

	appConfig.NextTryConnTime=0

	err:=autoSyncConfigServices(nil)
	err=autoSyncConfigServices(nil)

	test.Nil(t,err)

	config:=GetCurrentApolloConfig()

	test.Equal(t,"100004458",config.AppId)
	test.Equal(t,"default",config.Cluster)
	test.Equal(t,"application",config.NamespaceName)
	test.Equal(t,"20170430092936-dee2d58e74515ff3",config.ReleaseKey)
	//test.Equal(t,"value1",config.Configurations["key1"])
	//test.Equal(t,"value2",config.Configurations["key2"])
}

func TestAutoSyncConfigServicesNormal2NotModified(t *testing.T) {
	go runMockConfigServer(longNotmodifiedConfigResponse)
	defer closeMockConfigServer()

	time.Sleep(1*time.Second)

	appConfig.NextTryConnTime=0

	autoSyncConfigServicesSuccessCallBack([]byte(configResponseStr))

	config:=GetCurrentApolloConfig()

	fmt.Println("sleeping 10s")

	time.Sleep(10*time.Second)

	fmt.Println("checking cache time left")
	it := apolloConfigCache.NewIterator()
	for i := int64(0); i < apolloConfigCache.EntryCount(); i++ {
		entry := it.Next()
		if entry==nil{
			break
		}
		timeLeft, err := apolloConfigCache.TTL([]byte(entry.Key))
		test.Nil(t,err)
		fmt.Printf("key:%s,time:%v \n",string(entry.Key),timeLeft)
		test.Equal(t,timeLeft>=110,true)
	}

	test.Equal(t,"100004458",config.AppId)
	test.Equal(t,"default",config.Cluster)
	test.Equal(t,"application",config.NamespaceName)
	test.Equal(t,"20170430092936-dee2d58e74515ff3",config.ReleaseKey)
	test.Equal(t,"value1",getValue("key1"))
	test.Equal(t,"value2",getValue("key2"))

	err:=autoSyncConfigServices(nil)

	fmt.Println("checking cache time left")
	it1 := apolloConfigCache.NewIterator()
	for i := int64(0); i < apolloConfigCache.EntryCount(); i++ {
		entry := it1.Next()
		if entry==nil{
			break
		}
		timeLeft, err := apolloConfigCache.TTL([]byte(entry.Key))
		test.Nil(t,err)
		fmt.Printf("key:%s,time:%v \n",string(entry.Key),timeLeft)
		test.Equal(t,timeLeft>=120,true)
	}

	fmt.Println(err)
}

//test if not modify
func TestAutoSyncConfigServicesNotModify(t *testing.T) {
	go runMockConfigServer(notModifyConfigResponse)
	defer closeMockConfigServer()

	apolloConfig,err:=createApolloConfigWithJson([]byte(configResponseStr))
	updateApolloConfig(apolloConfig)

	time.Sleep(10*time.Second)
	checkCacheLeft(t,configCacheExpireTime-10)

	appConfig.NextTryConnTime=0

	err=autoSyncConfigServices(nil)

	test.Nil(t,err)

	config:=GetCurrentApolloConfig()

	test.Equal(t,"100004458",config.AppId)
	test.Equal(t,"default",config.Cluster)
	test.Equal(t,"application",config.NamespaceName)
	test.Equal(t,"20170430092936-dee2d58e74515ff3",config.ReleaseKey)

	checkCacheLeft(t,configCacheExpireTime)

	//test.Equal(t,"value1",config.Configurations["key1"])
	//test.Equal(t,"value2",config.Configurations["key2"])
}



func TestAutoSyncConfigServicesError(t *testing.T) {
	//reload app properties
	go initFileConfig()
	go runMockConfigServer(errorConfigResponse)
	defer closeMockConfigServer()

	time.Sleep(1*time.Second)

	err:=autoSyncConfigServices(nil)

	test.NotNil(t,err)

	config:=GetCurrentApolloConfig()

	//still properties config
	test.Equal(t,"test",config.AppId)
	test.Equal(t,"dev",config.Cluster)
	test.Equal(t,"application",config.NamespaceName)
	test.Equal(t,"",config.ReleaseKey)
}