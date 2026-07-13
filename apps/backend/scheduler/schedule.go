package scheduler

import (
	"log"
	"time"

	"resuming/database"
)

func DatabaseSchedulerGroupToSessionStore() {
	if err := database.SyncGroupErrorLogSessionStore(); err != nil {
		log.Println(err)
	}
	if err := database.SyncGroupClientAuditLogSessionStore(); err != nil {
		log.Println(err)
	}
	if err := database.SyncGroupAdminAuditLogSessionStore(); err != nil {
		log.Println(err)
	}
	if err := database.SyncGroupClientsConfigSessionStore(); err != nil {
		log.Println(err)
	}
	if err := database.SyncGroupAdminsConfigSessionStore(); err != nil {
		log.Println(err)
	}
	if err := database.SyncGroupClientReportLogSessionStore(); err != nil {
		log.Println(err)
	}
	if err := database.SyncGroupClientSupportMessagingSessionStore(); err != nil {
		log.Println(err)
	}
	if err := database.SyncGroupShowCaseRecordsSessionStore(); err != nil {
		log.Println(err)
	}
	if err := database.SyncGroupPortfoliosSessionStore(); err != nil {
		log.Println(err)
	}
	if err := database.SyncGroupResumesSessionStore(); err != nil {
		log.Println(err)
	}
	if err := database.SyncGroupAtsSessionStore(); err != nil {
		log.Println(err)
	}
	if err := database.SyncGroupUsersSessionStore(); err != nil {
		log.Println(err)
	}
	if err := database.SyncGroupSessionsSessionStore(); err != nil {
		log.Println(err)
	}
}

func DatabaseSchedulerGroupToDatabase() {
	if err := database.SyncGroupErrorLogDatabase(); err != nil {
		log.Println(err)
	}
	if err := database.SyncGroupClientAuditLogDatabase(); err != nil {
		log.Println(err)
	}
	if err := database.SyncGroupAdminAuditLogDatabase(); err != nil {
		log.Println(err)
	}
	if err := database.SyncGroupClientsConfigDatabase(); err != nil {
		log.Println(err)
	}
	if err := database.SyncGroupAdminsConfigDatabase(); err != nil {
		log.Println(err)
	}
	if err := database.SyncGroupClientReportLogDatabase(); err != nil {
		log.Println(err)
	}
	if err := database.SyncGroupClientSupportMessagingDatabase(); err != nil {
		log.Println(err)
	}
	if err := database.SyncGroupShowCaseRecordsDatabase(); err != nil {
		log.Println(err)
	}
	if err := database.SyncGroupPortfoliosDatabase(); err != nil {
		log.Println(err)
	}
	if err := database.SyncGroupResumesDatabase(); err != nil {
		log.Println(err)
	}
	if err := database.SyncGroupAtsDatabase(); err != nil {
		log.Println(err)
	}
	if err := database.SyncGroupUsersDatabase(); err != nil {
		log.Println(err)
	}
	if err := database.SyncGroupSessionsDatabase(); err != nil {
		log.Println(err)
	}
}

func FirstSync() {
	DatabaseSchedulerGroupToSessionStore()
}

func FullSync() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		DatabaseSchedulerGroupToDatabase()
		DatabaseSchedulerGroupToSessionStore()
	}
}
