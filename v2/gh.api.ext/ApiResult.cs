using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;

namespace gh.api.ext
{
    public class ApiResult
    {
        public object Result { get; set; }
        public DateTimeOffset CreatedAt { get; set; }
    }
}
