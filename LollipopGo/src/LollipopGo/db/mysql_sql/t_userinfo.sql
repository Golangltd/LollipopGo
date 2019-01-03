/*
Navicat MySQL Data Transfer

Source Server         : www.xshooting.com
Source Server Version : 50173
Source Host           : 47.107.125.75:3306
Source Database       : gl_XiaoMq

Target Server Type    : MYSQL
Target Server Version : 50173
File Encoding         : 65001

Date: 2019-01-03 16:35:05
*/

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for `t_userinfo`
-- ----------------------------
DROP TABLE IF EXISTS `t_userinfo`;
CREATE TABLE `t_userinfo` (
  `ID` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `uid` int(11) NOT NULL,
  `openid` varchar(70) NOT NULL,
  `vip` int(20) NOT NULL,
  `name` varchar(20) NOT NULL,
  `headurl` varchar(200) NOT NULL,
  `school` varchar(30) NOT NULL,
  `sex` varchar(30) NOT NULL,
  `hallexp` int(30) NOT NULL,
  `coinnum` int(30) NOT NULL,
  `masonrynum` int(30) NOT NULL,
  `mcard` int(30) NOT NULL,
  `constellation` varchar(30) NOT NULL,
  `medallist` varchar(200) NOT NULL,
  `createtime` varchar(20) NOT NULL,
  PRIMARY KEY (`ID`)
) ENGINE=MyISAM AUTO_INCREMENT=68 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of t_userinfo
-- ----------------------------
INSERT INTO `t_userinfo` VALUES ('67', '787', '3621f1454cacf995530ea53652ddf8fb', '0', '德玛西亚', 'http://xmqvip1-1253933147.file.myqcloud.com/ugc/images/2018/04/16/152387637868832h.jpg', '四川大学锦城学院', '1', '0', '2000', '0', '0', '金牛座', '', '1546503672');
