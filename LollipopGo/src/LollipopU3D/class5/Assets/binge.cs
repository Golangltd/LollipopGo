using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using System;
using Warfare.Coding;

public class binge : MonoBehaviour {

    // 程序启动时候执行
    void Start()
    {
        //GET请求 登录服务器获取数据
        StartCoroutine(GET("http://local.golang.ltd:8891/GolangLtdDT?Protocol=8&Protocol2=1&LoginName=binge&LoginPW=wswss&Timestamp=123377488"));
    }

    // 每帧执行一次
    void Update()
    {
    }

    //GET请求  
    IEnumerator GET(string url)
    {

        WWW www = new WWW(url);
        yield return www;

        if (www.error != null)
        {
            //GET请求失败   
            Debug.Log("error is :" + www.error);

        }
        else
        {
            //GET请求成功   
            Debug.Log("request ok : " + www.text);
            // 解析base64数据
            string decodeMessage = UnicodeConverter.base64decode(www.text);
            string utf8Message = UnicodeConverter.utf8to16(decodeMessage);
            Debug.Log("data : " + utf8Message);
            // json 数据解析  -->   ws 转
            // 数据 websocket操作
            // 获取网关服务器地址，建立长链接！
            // 启动心跳测试
            // 主要测试消息流程等
        }
    }

    public void Binge_http()
    {
        Debug.Log("点击了按钮！！！");
    }
}
