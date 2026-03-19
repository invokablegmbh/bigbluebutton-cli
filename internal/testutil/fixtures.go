package testutil

// fixtures maps API call names to canned XML/JSON responses.
var fixtures = map[string]string{
	"create": `<response>
	<returncode>SUCCESS</returncode>
	<meetingID>test-meeting-1</meetingID>
	<internalMeetingID>abc123def456</internalMeetingID>
	<parentMeetingID>bbb-none</parentMeetingID>
	<attendeePW>ap</attendeePW>
	<moderatorPW>mp</moderatorPW>
	<createTime>1699900000000</createTime>
	<voiceBridge>70001</voiceBridge>
	<dialNumber>613-555-1234</dialNumber>
	<createDate>Mon Nov 13 12:00:00 UTC 2023</createDate>
	<hasUserJoined>false</hasUserJoined>
	<duration>0</duration>
	<hasBeenForciblyEnded>false</hasBeenForciblyEnded>
</response>`,

	"join": `<response>
	<returncode>SUCCESS</returncode>
	<meeting_id>test-meeting-1</meeting_id>
	<user_id>user-123</user_id>
	<auth_token>auth-abc</auth_token>
	<session_token>sess-xyz</session_token>
	<url>https://bbb.example.com/client/BigBlueButton.html?sessionToken=sess-xyz</url>
</response>`,

	"end": `<response>
	<returncode>SUCCESS</returncode>
	<messageKey>sentEndMeetingRequest</messageKey>
	<message>A]request to end the meeting was sent. Please wait a few seconds, and then use the getMeetingInfo or isMeetingRunning API calls to verify that it was ended.</message>
</response>`,

	"insertDocument": `<response>
	<returncode>SUCCESS</returncode>
</response>`,

	"isMeetingRunning": `<response>
	<returncode>SUCCESS</returncode>
	<running>true</running>
</response>`,

	"getMeetingInfo": `<response>
	<returncode>SUCCESS</returncode>
	<meetingName>Test Meeting</meetingName>
	<meetingID>test-meeting-1</meetingID>
	<internalMeetingID>abc123def456</internalMeetingID>
	<createTime>1699900000000</createTime>
	<createDate>Mon Nov 13 12:00:00 UTC 2023</createDate>
	<voiceBridge>70001</voiceBridge>
	<dialNumber>613-555-1234</dialNumber>
	<attendeePW>ap</attendeePW>
	<moderatorPW>mp</moderatorPW>
	<running>true</running>
	<duration>0</duration>
	<hasUserJoined>true</hasUserJoined>
	<recording>false</recording>
	<hasBeenForciblyEnded>false</hasBeenForciblyEnded>
	<startTime>1699900000000</startTime>
	<endTime>0</endTime>
	<participantCount>2</participantCount>
	<listenerCount>1</listenerCount>
	<voiceParticipantCount>1</voiceParticipantCount>
	<videoCount>1</videoCount>
	<maxUsers>0</maxUsers>
	<moderatorCount>1</moderatorCount>
	<attendees>
		<attendee>
			<userID>user-1</userID>
			<fullName>John Doe</fullName>
			<role>MODERATOR</role>
			<isPresenter>true</isPresenter>
			<isListeningOnly>false</isListeningOnly>
			<hasJoinedVoice>true</hasJoinedVoice>
			<hasVideo>true</hasVideo>
			<clientType>HTML5</clientType>
		</attendee>
		<attendee>
			<userID>user-2</userID>
			<fullName>Jane Smith</fullName>
			<role>VIEWER</role>
			<isPresenter>false</isPresenter>
			<isListeningOnly>true</isListeningOnly>
			<hasJoinedVoice>false</hasJoinedVoice>
			<hasVideo>false</hasVideo>
			<clientType>HTML5</clientType>
		</attendee>
	</attendees>
	<isBreakout>false</isBreakout>
</response>`,

	"getMeetings": `<response>
	<returncode>SUCCESS</returncode>
	<meetings>
		<meeting>
			<meetingName>Test Meeting 1</meetingName>
			<meetingID>meeting-1</meetingID>
			<internalMeetingID>int-meeting-1</internalMeetingID>
			<createTime>1699900000000</createTime>
			<createDate>Mon Nov 13 12:00:00 UTC 2023</createDate>
			<voiceBridge>70001</voiceBridge>
			<dialNumber>613-555-1234</dialNumber>
			<attendeePW>ap</attendeePW>
			<moderatorPW>mp</moderatorPW>
			<running>true</running>
			<duration>0</duration>
			<hasUserJoined>true</hasUserJoined>
			<recording>false</recording>
			<hasBeenForciblyEnded>false</hasBeenForciblyEnded>
			<startTime>1699900000000</startTime>
			<endTime>0</endTime>
			<participantCount>5</participantCount>
			<listenerCount>3</listenerCount>
			<voiceParticipantCount>2</voiceParticipantCount>
			<videoCount>2</videoCount>
			<maxUsers>0</maxUsers>
			<moderatorCount>1</moderatorCount>
			<isBreakout>false</isBreakout>
		</meeting>
		<meeting>
			<meetingName>Test Meeting 2</meetingName>
			<meetingID>meeting-2</meetingID>
			<internalMeetingID>int-meeting-2</internalMeetingID>
			<createTime>1699910000000</createTime>
			<createDate>Mon Nov 13 15:00:00 UTC 2023</createDate>
			<voiceBridge>70002</voiceBridge>
			<dialNumber>613-555-5678</dialNumber>
			<attendeePW>ap2</attendeePW>
			<moderatorPW>mp2</moderatorPW>
			<running>true</running>
			<duration>60</duration>
			<hasUserJoined>true</hasUserJoined>
			<recording>true</recording>
			<hasBeenForciblyEnded>false</hasBeenForciblyEnded>
			<startTime>1699910000000</startTime>
			<endTime>0</endTime>
			<participantCount>10</participantCount>
			<listenerCount>8</listenerCount>
			<voiceParticipantCount>5</voiceParticipantCount>
			<videoCount>3</videoCount>
			<maxUsers>20</maxUsers>
			<moderatorCount>2</moderatorCount>
			<isBreakout>false</isBreakout>
		</meeting>
	</meetings>
</response>`,

	"getRecordings": `<response>
	<returncode>SUCCESS</returncode>
	<recordings>
		<recording>
			<recordID>rec-1</recordID>
			<meetingID>meeting-1</meetingID>
			<internalMeetingID>int-meeting-1</internalMeetingID>
			<name>Test Recording 1</name>
			<isBreakout>false</isBreakout>
			<published>true</published>
			<state>published</state>
			<startTime>1699900000000</startTime>
			<endTime>1699903600000</endTime>
			<participants>5</participants>
			<playback>
				<format>
					<type>presentation</type>
					<url>https://bbb.example.com/playback/presentation/2.3/rec-1</url>
					<processingTime>120000</processingTime>
					<length>60</length>
				</format>
			</playback>
		</recording>
		<recording>
			<recordID>rec-2</recordID>
			<meetingID>meeting-2</meetingID>
			<internalMeetingID>int-meeting-2</internalMeetingID>
			<name>Test Recording 2</name>
			<isBreakout>false</isBreakout>
			<published>false</published>
			<state>unpublished</state>
			<startTime>1699910000000</startTime>
			<endTime>1699913600000</endTime>
			<participants>10</participants>
			<playback>
				<format>
					<type>presentation</type>
					<url>https://bbb.example.com/playback/presentation/2.3/rec-2</url>
					<processingTime>180000</processingTime>
					<length>90</length>
				</format>
			</playback>
		</recording>
	</recordings>
</response>`,

	"publishRecordings": `<response>
	<returncode>SUCCESS</returncode>
	<published>true</published>
</response>`,

	"deleteRecordings": `<response>
	<returncode>SUCCESS</returncode>
	<deleted>true</deleted>
</response>`,

	"updateRecordings": `<response>
	<returncode>SUCCESS</returncode>
	<updated>true</updated>
</response>`,

	"getRecordingTextTracks": `{"response":{"returncode":"SUCCESS","tracks":[{"href":"https://bbb.example.com/textTrack/rec-1/en.vtt","kind":"subtitles","label":"English","lang":"en","source":"upload"}]}}`,

	"putRecordingTextTrack": `{"response":{"returncode":"SUCCESS","messageKey":"upload_text_track_success","message":"Text track uploaded successfully"}}`,

	"sendChatMessage": `<response>
	<returncode>SUCCESS</returncode>
	<messageKey>sentChatMessage</messageKey>
	<message>Chat message sent successfully</message>
</response>`,

	"getJoinUrl": `<response>
	<returncode>SUCCESS</returncode>
	<url>https://bbb.example.com/bigbluebutton/api/join?sessionToken=sess-xyz</url>
</response>`,
}
