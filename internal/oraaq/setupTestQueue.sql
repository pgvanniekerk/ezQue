-- Create a type for the AQ message that aligns with Oracle JMS Text Message.
-- Important to be aligned with dequeueSQL block.
BEGIN
    DBMS_AQADM.CREATE_QUEUE_TABLE(
        queue_table            =>  'text_msg_queue_table',
        queue_payload_type     =>  'SYS.AQ$_JMS_TEXT_MESSAGE',
        compatible             =>  '8.1',
        storage_clause         =>  'TABLESPACE USERS');
END;
/

-- Create the Advanced Queue on the Queue Table just created.
BEGIN
    DBMS_AQADM.CREATE_QUEUE(
        queue_name     =>  'text_msg_queue',
        queue_table    =>  'text_msg_queue_table');
END;
/

-- Then start the created queue.
BEGIN
    DBMS_AQADM.START_QUEUE('text_msg_queue');
END;
/