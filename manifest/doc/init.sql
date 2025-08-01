CREATE DATABASE IF NOT EXISTS `financial` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE `financial`;

SET NAMES utf8mb4;

DROP TABLE IF EXISTS `index_sample`;
CREATE TABLE `index_sample` (
    `type_code` VARCHAR(20) NOT NULL COMMENT '类型代码（中证指数，www.csindex.com.cn）',
    `type_name` VARCHAR(20) NOT NULL COMMENT '类型名称（沪深300、中证500、上证50....）',
    `stock_code` CHAR(6) NOT NULL COMMENT '股票代码'
) ENGINE=InnoDB DEFAULT CHARSET=UTF8 COMMENT='指数样本（来源：中证指数）';

DROP TABLE IF EXISTS `category`;
CREATE TABLE `category` (
    `type` VARCHAR(5) NOT NULL COMMENT '分类类型（证券会、中证）',
    `code` VARCHAR(10) NOT NULL COMMENT '行业Code',
    `name` VARCHAR(20) NOT NULL COMMENT '名称',
    `level` VARCHAR(2) NOT NULL COMMENT '层级',
    `display_order` TINYINT(3) UNSIGNED DEFAULT NULL COMMENT '显示顺序',
    `parent_code` CHAR(8) DEFAULT NULL COMMENT '父分类Code',
    PRIMARY KEY (`code`),
    KEY `i_parent_code` (`parent_code`) USING BTREE,
    CONSTRAINT `fk_parent_code` FOREIGN KEY (`parent_code`) REFERENCES `category` (`code`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=UTF8 COMMENT='行业分类';

DROP TABLE IF EXISTS `stock`;
CREATE TABLE `stock` (
    `code` CHAR(6) NOT NULL COMMENT '代码',
    `name` VARCHAR(16) DEFAULT NULL COMMENT '名称',
    `name_pinyin` VARCHAR(16) DEFAULT NULL COMMENT '名称（拼音）',
    `before_name` VARCHAR(100) DEFAULT NULL COMMENT '曾用名称',
    `company_name` VARCHAR(50) DEFAULT NULL COMMENT '公司名称',
    `company_type` VARCHAR(10) DEFAULT NULL COMMENT '公司类型',
    `company_type_code` CHAR(1) DEFAULT NULL COMMENT '公司类型代码',
    `company_profile` TEXT DEFAULT NULL COMMENT '公司简介',
    `region` VARCHAR(8) DEFAULT NULL COMMENT '地域（省份）',
    `address` TEXT DEFAULT NULL COMMENT '办公地址',
    `website` VARCHAR(50) DEFAULT NULL COMMENT '公司网站',
    `main_business` TEXT DEFAULT NULL COMMENT '主营业务',
    `business_scope` TEXT DEFAULT NULL COMMENT '经营范围',
    `create_date` DATE DEFAULT NULL COMMENT '成立日期',
    `listing_date` DATE DEFAULT NULL COMMENT '上市日期',
    `law_firm` VARCHAR(100) DEFAULT NULL COMMENT '律师事务所',
    `accounting_firm` VARCHAR(50) DEFAULT NULL COMMENT '会计师事务所',
    `market_place` CHAR(2) DEFAULT NULL COMMENT '交易市场（上海、深圳、北京）',
    PRIMARY KEY (`code`),
    KEY `i_name` (`name`),
    KEY `i_name_pinyin` (`name_pinyin`)
) ENGINE=InnoDB DEFAULT CHARSET=UTF8 COMMENT='股票';