-- First stop the queue.
BEGIN
    DBMS_AQADM.STOP_QUEUE('text_msg_queue');
END;
/

-- Then drop the queue.
BEGIN
    DBMS_AQADM.DROP_QUEUE('text_msg_queue');
END;
/

-- Finally drop the queue table.
BEGIN
    DBMS_AQADM.DROP_QUEUE_TABLE('text_msg_queue_table', TRUE, FALSE);
END;
/