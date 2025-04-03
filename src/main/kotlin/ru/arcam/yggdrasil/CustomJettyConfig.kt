package ru.arcam.yggdrasil

import org.slf4j.LoggerFactory
import org.slf4j.Logger
import org.springframework.boot.web.embedded.jetty.JettyServletWebServerFactory
import org.springframework.boot.web.server.WebServerFactoryCustomizer
import org.springframework.context.annotation.Bean
import org.springframework.context.annotation.Configuration


@Configuration
class CustomJettyConfig : WebServerFactoryCustomizer<JettyServletWebServerFactory> {
    val logger: Logger = LoggerFactory.getLogger(CustomJettyConfig::class.java)
    override fun customize(factory: JettyServletWebServerFactory) {
        logger.info("Configuring Tomcat...");
        factory.port = 8080
    }
    @Bean
    fun jettyFactory(): JettyServletWebServerFactory {
        return JettyServletWebServerFactory(8080)
    }
}