#encoding: utf-8

import json
from datetime import datetime

class DateTimeJSONEncoder(json.JSONEncoder):
    def default(self, o):
        if isinstance(o, datetime):
            return o.strftime('%Y-%m-%d %H:%M:%S')
        return super(DateTimeJSONEncoder, self).default(o)


if __name__ == '__main__':
    j = {
        'name' : 'kk',
        'age' : 29,
        'birthday' : datetime.strptime('1988-10-19', '%Y-%m-%d'),
    }
    print(json.dumps(j, cls=DateTimeJSONEncoder))
