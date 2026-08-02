#include "vita_bridge.h"

#if defined(PSP2) || defined(__Vita__)
#include <psp2/touch.h>
#include <psp2/motion.h>
#include <psp2/power.h>
#include <psp2/ctrl.h>

void vita_bridge_init(void) {
    sceTouchSetSamplingState(SCE_TOUCH_PORT_FRONT, SCE_TOUCH_SAMPLING_STATE_START);
    sceTouchSetSamplingState(SCE_TOUCH_PORT_BACK, SCE_TOUCH_SAMPLING_STATE_START);
    sceMotionStartSampling();
    sceCtrlSetSamplingMode(SCE_CTRL_MODE_DIGITALANALOG);
}

int vita_bridge_get_touch_front(VitaTouchPoint* out_points, int max_count) {
    SceTouchData touch_data;
    sceTouchPeek(SCE_TOUCH_PORT_FRONT, &touch_data, 1);
    int count = touch_data.reportNum < max_count ? touch_data.reportNum : max_count;
    for (int i = 0; i < count; i++) {
        out_points[i].id = touch_data.report[i].id;
        out_points[i].x = (touch_data.report[i].x * 960) / 1920;
        out_points[i].y = (touch_data.report[i].y * 544) / 1088;
        out_points[i].force = (float)touch_data.report[i].force / 255.0f;
    }
    return count;
}

int vita_bridge_get_touch_back(VitaTouchPoint* out_points, int max_count) {
    SceTouchData touch_data;
    sceTouchPeek(SCE_TOUCH_PORT_BACK, &touch_data, 1);
    int count = touch_data.reportNum < max_count ? touch_data.reportNum : max_count;
    for (int i = 0; i < count; i++) {
        out_points[i].id = touch_data.report[i].id;
        out_points[i].x = (touch_data.report[i].x * 960) / 1920;
        out_points[i].y = (touch_data.report[i].y * 544) / 1088;
        out_points[i].force = (float)touch_data.report[i].force / 255.0f;
    }
    return count;
}

void vita_bridge_get_gyro(VitaVector3* out_gyro) {
    SceMotionState motion_state;
    sceMotionGetState(&motion_state);
    out_gyro->x = motion_state.angularVelocity.x;
    out_gyro->y = motion_state.angularVelocity.y;
    out_gyro->z = motion_state.angularVelocity.z;
}

void vita_bridge_get_accel(VitaVector3* out_accel) {
    SceMotionState motion_state;
    sceMotionGetState(&motion_state);
    out_accel->x = motion_state.acceleration.x;
    out_accel->y = motion_state.acceleration.y;
    out_accel->z = motion_state.acceleration.z;
}

int vita_bridge_get_battery_level(void) {
    return scePowerGetBatteryLifePercent();
}

int vita_bridge_is_charging(void) {
    return scePowerIsBatteryCharging();
}

int vita_bridge_set_cpu_clock(int freq_mhz) {
    return scePowerSetArmClockFrequency(freq_mhz);
}

int main(int argc, char *argv[]) {
    (void)argc;
    (void)argv;
    vita_bridge_init();
    // WASM から変換された C モジュールの初期化
    return 0;
}

#else

void vita_bridge_init(void) {}
int vita_bridge_get_touch_front(VitaTouchPoint* out_points, int max_count) { (void)out_points; (void)max_count; return 0; }
int vita_bridge_get_touch_back(VitaTouchPoint* out_points, int max_count) { (void)out_points; (void)max_count; return 0; }
void vita_bridge_get_gyro(VitaVector3* out_gyro) { if(out_gyro) { out_gyro->x = 0; out_gyro->y = 0; out_gyro->z = 0; } }
void vita_bridge_get_accel(VitaVector3* out_accel) { if(out_accel) { out_accel->x = 0; out_accel->y = 0; out_accel->z = -1.0f; } }
int vita_bridge_get_battery_level(void) { return 100; }
int vita_bridge_is_charging(void) { return 1; }
int vita_bridge_set_cpu_clock(int freq_mhz) { (void)freq_mhz; return 0; }

int main(int argc, char *argv[]) {
    (void)argc;
    (void)argv;
    return 0;
}

#endif


