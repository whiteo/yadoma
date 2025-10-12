package dev.whiteo.yadoma.config;

import org.springframework.context.annotation.Configuration;
import org.springframework.data.mongodb.config.EnableMongoAuditing;

/**
 * Configuration class for enabling MongoDB auditing.
 * Auditing allows automatic population of creation and modification dates in entities.
 *
 * @author Leo Tanas (<a href="https://github.com/whiteo">github</a>)
 */
@Configuration
@EnableMongoAuditing
public class MongoConfig {}