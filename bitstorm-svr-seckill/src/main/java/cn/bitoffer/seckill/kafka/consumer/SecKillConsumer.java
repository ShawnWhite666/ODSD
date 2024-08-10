package cn.bitoffer.seckill.kafka.consumer;

import cn.bitoffer.seckill.service.SecKillService;
import cn.bitoffer.seckill.service.impl.SecKillServiceImpl;
import com.google.gson.Gson;
import lombok.extern.slf4j.Slf4j;
import org.apache.kafka.clients.consumer.ConsumerRecord;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.context.annotation.Bean;
import org.springframework.kafka.annotation.KafkaListener;
import org.springframework.kafka.config.ConcurrentKafkaListenerContainerFactory;
import org.springframework.kafka.config.KafkaListenerContainerFactory;
import org.springframework.kafka.core.ConsumerFactory;
import org.springframework.kafka.listener.AbstractMessageListenerContainer;
import org.springframework.kafka.listener.ConcurrentMessageListenerContainer;
import org.springframework.kafka.support.Acknowledgment;
import org.springframework.kafka.support.KafkaHeaders;
import org.springframework.messaging.handler.annotation.Header;
import org.springframework.stereotype.Component;

import java.util.Optional;

@Component
@Slf4j
public class SecKillConsumer {
    @Autowired
    private SecKillService secKillService;
    @KafkaListener(topics = "tp-seckill", groupId = "TEST_GROUP",concurrency = "1", containerFactory = "kafkaManualAckListenerContainerFactory")
    public void topic_test(ConsumerRecord<?, ?> record, Acknowledgment ack, @Header(KafkaHeaders.RECEIVED_TOPIC) String topic) {
        Optional message = Optional.ofNullable(record.value());
        if (message.isPresent()) {
            Object msg = message.get();
            System.out.println("aaaaaa收到Kafka消息! Topic:" + topic + ",Message:" + msg);
            try {
                // 这里写你对接收到的消息的处理逻辑
                // 手动ACK
                Gson gson = new Gson();
                SecKillServiceImpl.SecKillMsg skMsg = new SecKillServiceImpl.SecKillMsg();
                skMsg = gson.fromJson(msg.toString(), SecKillServiceImpl.SecKillMsg.class);
                secKillService.secKillInStore(skMsg.traceID, skMsg.userID, skMsg.goods,  skMsg.secNum, skMsg.num);
                ack.acknowledge();
                log.info("Kafka消费成功! Topic:" + topic + ",Message:" + msg);
            } catch (Exception e) {
                e.printStackTrace();
                log.error("Kafka消费失败！Topic:" + topic + ",Message:" + msg, e);
            }
        }
    }

}




