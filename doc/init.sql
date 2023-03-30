CREATE DATABASE IF NOT EXISTS `financial` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE `financial`;

SET NAMES utf8mb4;

DROP TABLE IF EXISTS `index`;
CREATE TABLE `index_sample` (
    `type_code` VARCHAR(20) NOT NULL COMMENT '类型代码（中证指数，www.csindex.com.cn）',
    `type_name` VARCHAR(20) NOT NULL COMMENT '类型名称（沪深300、中证500、上证50....）',
    `stock_code` CHAR(6) NOT NULL COMMENT '股票代码'
) ENGINE=InnoDB DEFAULT CHARSET=UTF8 COMMENT='指数样本信息（来源：中证指数）';

DROP TABLE IF EXISTS `category`;
CREATE TABLE `category` (
    `type` VARCHAR(5) NOT NULL COMMENT '分类类型（证券会、中证）',
    `id` VARCHAR(10) NOT NULL COMMENT '行业ID',
    `name` VARCHAR(20) NOT NULL COMMENT '名称',
    `level` VARCHAR(2) NOT NULL COMMENT '层级',
    `display_order` TINYINT(3) UNSIGNED DEFAULT NULL COMMENT '显示顺序',
    `parent_id` CHAR(8) DEFAULT NULL COMMENT '父分类ID',
    PRIMARY KEY (`id`),
    KEY `i_parent_id` (`parent_id`) USING BTREE,
    CONSTRAINT `fk_parent_id` FOREIGN KEY (`parent_id`) REFERENCES `category` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=UTF8 COMMENT='行业分类';

DROP TABLE IF EXISTS `category_stock_code`;
CREATE TABLE `category_stock_code` (
    `type` VARCHAR(5) NOT NULL COMMENT '分类类型（证券会、中证）',
    `category_id` VARCHAR(10) NOT NULL COMMENT '行业ID',
    `stock_code` CHAR(6) NOT NULL COMMENT '股票代码',
    PRIMARY KEY (`category_id`, `stock_code`),
    CONSTRAINT `fk_category_id` FOREIGN KEY (`category_id`) REFERENCES `category` (`id`) ON DELETE CASCADE,
    KEY `i_stock_code` (`stock_code`)
) ENGINE=InnoDB DEFAULT CHARSET=UTF8 COMMENT='股票行业分类';

DROP TABLE IF EXISTS `stock`;
CREATE TABLE `stock` (
    `code` CHAR(6) NOT NULL COMMENT '股票代码',
    `stock_name` VARCHAR(16) DEFAULT NULL COMMENT '股票名称',
    `stock_name_pinyin` VARCHAR(16) DEFAULT NULL COMMENT '股票名称（拼音）',
    `stock_before_name` VARCHAR(100) DEFAULT NULL COMMENT '股票曾用名称',
    `company_name` VARCHAR(50) DEFAULT NULL COMMENT '公司名称',
    `company_profile` TEXT DEFAULT NULL COMMENT '公司简介',
    `region` VARCHAR(8) DEFAULT NULL COMMENT '地域（省份）',
    `address` TEXT DEFAULT NULL COMMENT '办公地址',
    `website` VARCHAR(50) DEFAULT NULL COMMENT '公司网站',
    `main_business` TEXT DEFAULT NULL COMMENT '主营业务',
    `business_scope` TEXT DEFAULT NULL COMMENT '经营范围',
    `date_of_incorporation` DATE DEFAULT NULL COMMENT '成立日期',
    `listing_date` DATE DEFAULT NULL COMMENT '上市日期',
    `law_firm` VARCHAR(100) DEFAULT NULL COMMENT '律师事务所',
    `accounting_firm` VARCHAR(50) DEFAULT NULL COMMENT '会计师事务所',
    `market_place` CHAR(2) DEFAULT NULL COMMENT '交易市场（上海、深圳、北京）',
    PRIMARY KEY (`code`),
    KEY `i_stock_name` (`stock_name`),
    KEY `i_stock_name_pinyin` (`stock_name_pinyin`)
) ENGINE=InnoDB DEFAULT CHARSET=UTF8 COMMENT='股票信息';

DROP TABLE IF EXISTS `financial`;
CREATE TABLE `financial` (
    `code` CHAR(6) NOT NULL COMMENT '股票代码',
    `year` CHAR(4) NOT NULL COMMENT '年份',
    `report_date` DATE NOT NULL COMMENT '财报季期',
    `report_type` VARCHAR(2) NOT NULL COMMENT '季期类型（Q1~Q4，分别代表：一季报、半年报、三季报、年报；O，代表：其他）',
    `dividend` DOUBLE DEFAULT NULL COMMENT '年度分红总金额',
    `ocf` DOUBLE DEFAULT NULL COMMENT '营业活动现金流量',
    `cfi` DOUBLE DEFAULT NULL COMMENT '投资活动现金流量',
    `cff` DOUBLE DEFAULT NULL COMMENT '筹资活动现金流量',
    `assign_dividend_porfit` DOUBLE DEFAULT NULL COMMENT '分配股利、利润或偿付利息支付的现金',
    `acquisition_assets` DOUBLE DEFAULT NULL COMMENT '购建固定资产、无形资产和其他长期资产支付的现金',
    `np` DOUBLE DEFAULT NULL COMMENT '净利润',
    `oi` DOUBLE DEFAULT NULL COMMENT '营业收入',
    `coe` DOUBLE DEFAULT NULL COMMENT '营业成本',
    `coe_total` DOUBLE DEFAULT NULL COMMENT '营业总成本（含各种费用，销售费用、管理费用等）',
    `eps` DOUBLE DEFAULT NULL COMMENT '每股盈余|基本每股收益',
    `monetary_fund` DOUBLE DEFAULT NULL COMMENT '货币资金',
    `trade_finasset_notfvtpl` DOUBLE DEFAULT NULL COMMENT '交易性金融资产',
    `trade_finasset` DOUBLE DEFAULT NULL COMMENT '交易性金融资产（历史遗留）',
    `derive_finasset` DOUBLE DEFAULT NULL COMMENT '衍生金融资产',
    `fixed_asset` DOUBLE DEFAULT NULL COMMENT '固定资产',
    `cip` DOUBLE DEFAULT NULL COMMENT '在建工程',
    `ca_total` DOUBLE DEFAULT NULL COMMENT '流动资产总额',
    `nca_total` DOUBLE DEFAULT NULL COMMENT '非流动资产总额',
    `cl_total` DOUBLE DEFAULT NULL COMMENT '流动负债总额',
    `ncl_total` DOUBLE DEFAULT NULL COMMENT '非流动负债总额',
    `inventory` DOUBLE DEFAULT NULL COMMENT '存货',
    `inventory_liquidating` DOUBLE DEFAULT NULL COMMENT '存货减少额',
    `accounts_rece` DOUBLE DEFAULT NULL COMMENT '应收账款',
    `accounts_payable` DOUBLE DEFAULT NULL COMMENT '应付账款',
    `np_ratio` DOUBLE DEFAULT NULL COMMENT '净利率：净利润 / 营业收入',
    `dividend_ratio` DOUBLE DEFAULT NULL COMMENT '分红率：分红总金额 / 净利润',
    `oi_ratio` DOUBLE DEFAULT NULL COMMENT '营业毛利率：(营业收入 - 营业成本) / 营业收入',
    `operating_profit_ratio` DOUBLE DEFAULT NULL COMMENT '营业利益率|营业利润率：(营业收入 - 营业成本 - 营业费用) / 营业收入',
    `operating_safety_ratio` DOUBLE DEFAULT NULL COMMENT '经营安全边际率：营业利益率 / 营业毛利率',
    `cash_equivalent_ratio` DOUBLE DEFAULT NULL COMMENT '现金与约当现金比率：(货币资金 + 交易性金融资产 + 衍生金融资产) / (流动资产总额 + 非流动资产总额)',
    `ca_ratio` DOUBLE DEFAULT NULL COMMENT '流动资产比率：流动资产总额 / (流动资产总额 + 非流动资产总额)',
    `cl_ratio` DOUBLE DEFAULT NULL COMMENT '流动负债比率：流动负债总额 / (流动资产总额 + 非流动资产总额)',
    `ncl_ratio` DOUBLE DEFAULT NULL COMMENT '长期负债比率：非流动负债总额 / (流动资产总额 + 非流动资产总额)',
    `debt_ratio` DOUBLE DEFAULT NULL COMMENT '负债比率：(流动负债总额 + 非流动负债总额) / (流动资产总额 + 非流动资产总额)',
    `long_term_funds_ratio` DOUBLE DEFAULT NULL COMMENT '长期资金占不动产及设备比率：(非流动负债总额 + 股东权益) / (固定资产 + 在建工程)',
    `equity_ratio` DOUBLE DEFAULT NULL COMMENT '股东权益比率：100 - 负债比率',
    `inventory_ratio` DOUBLE DEFAULT NULL COMMENT '存货比率：存货 / (流动资产总额 + 非流动资产总额)',
    `accounts_rece_ratio` DOUBLE DEFAULT NULL COMMENT '应收账款比率：应收账款 / (流动资产总额 + 非流动资产总额)',
    `accounts_payable_ratio` DOUBLE DEFAULT NULL COMMENT '应付账款比率：应付账款 / (流动资产总额 + 非流动资产总额)',
    `current_ratio` DOUBLE DEFAULT NULL COMMENT '流动比率：流动资产 / 流动负债总额',
    `quick_ratio` DOUBLE DEFAULT NULL COMMENT '速动比率：(流动资产 - 存货) / 流动负债总额',
    `roe` DOUBLE DEFAULT NULL COMMENT '股东权益报酬率：净利润 / (流动资产总额 + 非流动资产总额 - 流动负债总额 - 非流动负债总额)',
    `roa` DOUBLE DEFAULT NULL COMMENT '总资产报酬率：净利润 / (流动资产总额 + 非流动资产总额)',
    `accounts_rece_turnover_ratio` DOUBLE DEFAULT NULL COMMENT '应收账款周转率（次）：营业收入 / 应收账款',
    `average_cash_receipt_days` DOUBLE DEFAULT NULL COMMENT '平均收现天数：360 / 应收账款周转率',
    `inventory_turnover_ratio` DOUBLE DEFAULT NULL COMMENT '存货周转率（次）：营业成本 / 存货',
    `average_sales_days` DOUBLE DEFAULT NULL COMMENT '平均销货天数：360 / 存货周转率',
    `immovables_turnover_ratio` DOUBLE DEFAULT NULL COMMENT '不动产及设备周转率（次）：营业收入 / (固定资产 + 在建工程)',
    `total_asset_turnover_ratio` DOUBLE DEFAULT NULL COMMENT '总资产周转率（次）：营业收入 / (流动资产总额 + 非流动资产总额)',
    `cash_flow_ratio` DOUBLE DEFAULT NULL COMMENT '现金流量比率：营业活动现金流量 / 流动负债总额',
    `cash_flow_adequacy_ratio` DOUBLE DEFAULT NULL COMMENT '现金流量允当比率：近五年营业活动现金流量 / 近五年(购建固定资产、无形资产和其他长期资产支付的现金 + 分配股利、利润或偿付利息支付的现金 - 存货减少额)',
    `cash_reinvestment_ratio` DOUBLE DEFAULT NULL COMMENT '现金再投资比率：(经营活动产生的现金流量净额 - 现金股利) / (流动资产总额 + 非流动资产总额 - 流动负债总额)',
    `profit_cash_ratio` DOUBLE DEFAULT NULL COMMENT '盈利现金比率：营业活动现金流量 / 净利润',
    PRIMARY KEY (`code`, `year`, `report_date`)
) ENGINE=InnoDB DEFAULT CHARSET=UTF8 COMMENT='财务报表';
