install(
    TARGETS glitchfuck_exe
    RUNTIME COMPONENT glitchfuck_Runtime
)

if(PROJECT_IS_TOP_LEVEL)
  include(CPack)
endif()
