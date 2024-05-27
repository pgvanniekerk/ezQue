package oraaq

const dequeueSQL = `Declare
    queue_name          Varchar2(255) := :1;
    msgid               Raw(16);
    dequeue_options     DBMS_AQ.dequeue_options_t;
    message_properties  DBMS_AQ.message_properties_t;
    message             SYS.AQ$_JMS_TEXT_MESSAGE;
    extractedMessage    Clob;
                      
    errm                Varchar2(4000) := '';
            
Begin
    dequeue_options.dequeue_mode   := sys.DBMS_AQ.REMOVE;
    dequeue_options.wait           := DBMS_AQ.FOREVER;
    dequeue_options.visibility     := DBMS_AQ.ON_COMMIT;
    
    Begin
        DBMS_AQ.Dequeue(
            queue_name          => queue_name,
            dequeue_options     => dequeue_options,
            message_properties  => message_properties,
            payload             => message,
            msgid               => msgid
        );
        extractedMessage := message.text_vc;
    Exception
        When Others Then
            msgid := UTL_RAW.CAST_TO_RAW('');  -- Set to NULL explicitly
            extractedMessage := '';
            errm := SQLErrm;
    End;

    :2 := extractedMessage;
    :3 := RAWTOHEX(msgid);
    :4 := errm; -- no error

End;
`

const enqueueSql = `
	DECLARE
	   	enqueue_options     DBMS_AQ.ENQUEUE_OPTIONS_T;
	   	message_properties  DBMS_AQ.MESSAGE_PROPERTIES_T;
	   	message_handle      RAW(16);
	   	message             SYS.AQ$_JMS_TEXT_MESSAGE;

	    queue_name          Varchar2(255) := :1;
		msgContent 			Clob := :2;
	BEGIN
		message := SYS.AQ$_JMS_TEXT_MESSAGE.construct;
		message.set_text(msgContent);
		
		DBMS_AQ.ENQUEUE(
		  queue_name         => queue_name,
		  enqueue_options    => enqueue_options,
		  message_properties => message_properties,
		  payload            => message,
		  msgid              => message_handle
		);

		Commit;
	   
	END;
`
