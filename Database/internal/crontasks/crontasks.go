package crontasks

import (
	"context"
	"strconv"

	"digest_bot_database/internal/log"
	"digest_bot_database/internal/modules/sources"
)

func UpdateFullDigestsForUsers(ctx context.Context, sourcesModule *sources.Module) error {
	users, err := sourcesModule.Repository.GetUsersIDList(ctx)
	if err != nil {
		log.FromContext(ctx).Error(
			"get users ID list",
			"error", err,
		)
		return err
	}

	for _, user := range users {
		intUserID, err := strconv.Atoi(user)
		if err != nil {
			log.FromContext(ctx).Error(
				"atoi userID",
				"error", err,
			)
			continue
		}

		sourcesList, err := sourcesModule.Repository.GetSourcesByUserID(ctx, intUserID)
		if err != nil {
			log.FromContext(ctx).Error(
				"get sources for user by id",
				"error", err,
			)
			continue
		}

		for _, source := range sourcesList {
			newVids, err := sourcesModule.Client.GetNewVideosForUserSourceByHour(source)
			if err != nil {
				log.FromContext(ctx).Error(
					"get new videos for user source",
					"error", err,
				)
				continue
			}

			var newVidsIDs []string
			for _, vids := range newVids {
				newVidsIDs = append(newVidsIDs, vids.VideoID)
			}

			if err = sourcesModule.Repository.UpdateNewVidsForSource(ctx, newVidsIDs, source); err != nil {
				return err
			}
		}
	}
	return nil
}

func PurgeTodaysColums(ctx context.Context, sourcesModule *sources.Module) error {
	if err := sourcesModule.Repository.PurgeNewVidsAndTodaysDigestColumns(ctx); err != nil {
		log.FromContext(ctx).Error(
			"purge new vids and todays digest columns",
			"error", err,
		)
		return err
	}

	return nil
}
