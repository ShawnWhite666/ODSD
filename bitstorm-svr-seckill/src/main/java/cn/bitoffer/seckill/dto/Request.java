package cn.bitoffer.seckill.dto;

import lombok.Data;

import java.io.Serializable;
import java.util.Date;
import java.util.List;

/**
 * Post请求
 *
 */
@Data
public class ExamplePostRequest implements Serializable {

    /**
     * 名称
     */
    private String name;

    /**
     * 年龄
     */
    private int age;

    /**
     * 标签列表
     */
    private List<String> tags;

    private static final long serialVersionUID = 1L;
}

public class CreateTaskRequest {
    public TaskData taskData;
}

public class CreateTaskReply {
    public String taskId;
}

public class TaskData {
    public int    taskStatus;
    public long   priority;
    public long crtRetryNum;
    public long maxRetryNum;
    public long maxRetryInterval;
    public long orderTime;
    public String userId;
    public String taskId;
    public String taskType;
    public String taskStage;
    public String scheduleLog;
    public String context;
    public Date createTime;
    public Date modifyTime;


}