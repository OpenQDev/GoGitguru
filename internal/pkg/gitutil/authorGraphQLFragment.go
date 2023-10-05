package gitutil

var AUTHOR_GRAPHQL_FRAGMENT string = `fragment commitDetails on Commit {
	author {
			name
			email
			user {
					github_rest_id: databaseId
					github_graphql_id: id
					login
					name
					email
					avatar_url: avatarUrl
					company
					location
					hireable: isHireable
					bio
					blog: websiteUrl
					twitter_username: twitterUsername
					followers {
							totalCount
					}
					following {
							totalCount
					}
					created_at: createdAt
					updated_at: updatedAt
			}
	}
}`
