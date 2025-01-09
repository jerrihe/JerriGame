CREATE TABLE IF NOT EXISTS `JAccount` (
	`UserId` VARCHAR(128),
	`PlatType` VARCHAR(128),
	`AccId` BIGINT UNSIGNED,
	`CreateAt` BIGINT SIGNED,
	PRIMARY KEY(`AccId`),
	INDEX `UserId_PlatType_index` (`UserId`,`PlatType`)
) ENGINE=InnoDB DEFAULT CHARSET=UTF8MB4 COLLATE=utf8mb4_unicode_ci;