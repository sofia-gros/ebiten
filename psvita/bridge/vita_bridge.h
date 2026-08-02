#ifndef VITA_BRIDGE_H
#define VITA_BRIDGE_H

#include <stdint.h>

#ifdef __cplusplus
extern "C" {
#endif

typedef struct {
    int id;
    int x;
    int y;
    float force;
} VitaTouchPoint;

typedef struct {
    float x;
    float y;
    float z;
} VitaVector3;

// Bridge APIs
void vita_bridge_init(void);
int vita_bridge_get_touch_front(VitaTouchPoint* out_points, int max_count);
int vita_bridge_get_touch_back(VitaTouchPoint* out_points, int max_count);
void vita_bridge_get_gyro(VitaVector3* out_gyro);
void vita_bridge_get_accel(VitaVector3* out_accel);
int vita_bridge_get_battery_level(void);
int vita_bridge_is_charging(void);
int vita_bridge_set_cpu_clock(int freq_mhz);

#ifdef __cplusplus
}
#endif

#endif // VITA_BRIDGE_H
