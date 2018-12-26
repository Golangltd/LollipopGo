/*
Navicat MySQL Data Transfer

Source Server         : www.xshooting.com
Source Server Version : 50173
Source Host           : 47.107.125.75:3306
Source Database       : gl_XiaoMq

Target Server Type    : MYSQL
Target Server Version : 50173
File Encoding         : 65001

Date: 2018-12-26 17:30:59
*/

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for `t_userinfo_copy`
-- ----------------------------
DROP TABLE IF EXISTS `t_userinfo_copy`;
CREATE TABLE `t_userinfo_copy` (
  `ID` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `uid` int(11) NOT NULL,
  `vip` int(20) NOT NULL,
  `name` varchar(20) NOT NULL,
  `headurl` varchar(60) NOT NULL,
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
) ENGINE=MyISAM AUTO_INCREMENT=67 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of t_userinfo_copy
-- ----------------------------
INSERT INTO `t_userinfo_copy` VALUES ('66', '787', '0', '德玛西亚', 'http://xmqvip1-1253933147.file.myqcloud.com/ugc/images/2018/', '四川大学锦城学院', '1', '0', '2000', '8888', '1235', '金牛座', '', '1545801854');
