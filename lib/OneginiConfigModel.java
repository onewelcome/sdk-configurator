package com.onegini.example;

import android.os.Build;
import com.onegini.mobile.sdk.android.model.OneginiClientConfigModel;

public class OneginiConfigModel implements OneginiClientConfigModel {

  private final String appIdentifier = "value_will_be_replaced";
  private final String appPlatform = "android";
  private final String redirectionUri = "value_will_be_replaced";
  private final String appVersion = "value_will_be_replaced";
  private final String baseURL = "value_will_be_replaced";
  private final String resourceBaseURL = "value_will_be_replaced";
  private final String keystoreHash = "value_will_be_replaced";
  private final int maxPinFailures = 1;

  @Override
  public String getAppIdentifier() {
    return appIdentifier;
  }

  @Override
  public String getAppPlatform() {
    return appPlatform;
  }

  @Override
  public String getRedirectUri() {
    return redirectionUri;
  }

  @Override
  public String getAppVersion() {
    return appVersion;
  }

  @Override
  public String getBaseUrl() {
    return baseURL;
  }

  @Override
  public String getResourceBaseUrl() {
    return resourceBaseURL;
  }

  @Override
  public int getCertificatePinningKeyStore() {
    return R.raw.keystore;
  }

  @Override
  public String getKeyStoreHash() {
    return keystoreHash;
  }

  @Override
  public String getDeviceName() {
    return Build.BRAND + " " + Build.MODEL;
  }

  @Override
  public boolean shouldGetIdToken() {
    return false;
  }

  /**
   * Get the max PIN failures. This attribute is just used for visual representation towards the end-user.
   *
   * @Deprecated Since Android SDK 6.01.00 this attribute is fetched from the Device config.
   *
   * @return The max PIN failures
   */
  public int getMaxPinFailures() {
    return maxPinFailures;
  }

  @Override
  public String toString() {
    return "ConfigModel{" +
            "  appIdentifier='" + appIdentifier + "'" +
            ", appPlatform='" + appPlatform + "'" +
            ", redirectionUri='" + redirectionUri + "'" +
            ", appVersion='" + appVersion + "'" +
            ", baseURL='" + baseURL + "'" +
            ", maxPinFailures='" + maxPinFailures + "'" +
            ", resourceBaseURL='" + resourceBaseURL + "'" +
            ", keyStoreHash='" + getKeyStoreHash() + "'" +
            "}";
  }
}
