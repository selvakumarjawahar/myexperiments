
test_dict = {
    'connect_time_avg': 1,
    'first_byte_time_avg': 2,
    'response_time_min': 3,
    'response_time_avg': 4,
    'response_time_max': 5,
    'size_avg': 6,
    'speed_avg': 7,
    'es': 8,
    'es_timeout': 9,
    'es_response': 10
}
tasktype = 'http'
out_dict = {}

for key in test_dict:
    new_key = tasktype + '_' + key
    out_dict[new_key] = test_dict[key]
if 'es' in test_dict:
    out_dict['es'] = test_dict['es']

print out_dict

